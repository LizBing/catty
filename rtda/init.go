// Package rtda — shared class/interface initialization service (ADR-0025).
//
// This file implements the Java 25 single-execution-context initialization
// state machine used across Interpreter, IR, and AOT. Each runtime class or
// interface identity has exactly one of four states: not-initialized,
// initializing, initialized, or erroneous.
//
// Bytecode execution of <clinit> is provided by the caller via a callback,
// keeping rtda free of bytecode dispatch and avoiding import cycles.

package rtda

// InitResult encodes the outcome of a class/interface initialization request.
type InitResult struct {
	ErrObj *Object // non-nil on failure
}

// SuccessInit returns a successful initialization result.
func SuccessInit() InitResult { return InitResult{} }

// ClinitRunner executes a <clinit> method and returns its outcome.
// Implementations live in the interpreter and AOT bridge packages.
type ClinitRunner func(class *Class, clinit *Method) InitResult

// InitializeClass performs the JVMS §5.5 initialization procedure for one
// class or interface.
//
// Parameters:
//   - loader: for resolving predecessor classes
//   - class: the class/interface to initialize
//   - ecID: execution-context identity (for recursive-request recognition)
//   - runClinit: callback that executes <clinit> bytecode
func InitializeClass(loader Loader, class *Class, ecID uint64, runClinit ClinitRunner) InitResult {
	// Fast path: already initialized.
	if class.InitState() == initInitialized {
		return SuccessInit()
	}

	// Recursive same-owner request — return normally without re-running <clinit>.
	if class.InitState() == initInProgress && class.InitOwner() == ecID {
		return SuccessInit()
	}

	// Erroneous — the caller will throw NoClassDefFoundError.
	if class.InitState() == initErroneous {
		return buildNCDFE(loader, class.name)
	}

	// Claim initializing state BEFORE recursing into predecessors
	// (JVMS §5.5 step 6 precedes step 7). Same-owner recursive requests
	// that arrive during predecessor initialization will hit the
	// initInProgress+same-owner check above and return normally.
	if !class.MarkInitInProgress(ecID) {
		// Racing claim (won't happen in single-context R2). Re-check.
		return InitializeClass(loader, class, ecID, runClinit)
	}

	// For classes: recursively initialize superclass first, then
	// default-bearing superinterfaces in JVMS §5.5 step 7 order.
	if !class.IsInterface() {
		if class.superClass != nil {
			if r := InitializeClass(loader, class.superClass, ecID, runClinit); r.ErrObj != nil {
				class.MarkErroneous()
				return r
			}
		}
		for _, iface := range class.DefaultBearingSuperInterfaces(make(map[string]bool)) {
			if r := InitializeClass(loader, iface, ecID, runClinit); r.ErrObj != nil {
				class.MarkErroneous()
				return r
			}
		}
	}

	// Run <clinit> if present; otherwise mark initialized directly.
	clinit := class.GetMethod("<clinit>", "()V")
	if clinit != nil {
		if r := runClinit(class, clinit); r.ErrObj != nil {
			class.MarkErroneous()
			return r
		}
	}

	class.MarkInitialized()
	return SuccessInit()
}

// buildNCDFE creates a NoClassDefFoundError for an erroneous class.
func buildNCDFE(loader Loader, className string) InitResult {
	ncdfeClass := loader.LoadClass("java/lang/NoClassDefFoundError")
	ncdfe := NewObject(ncdfeClass)
	setDetailMessage(loader, ncdfe, className)
	return InitResult{ErrObj: ncdfe}
}

// wrapInitFailure wraps a thrown object per JVMS §5.5:
// If t is java.lang.Error or a subclass, propagate it directly.
// Otherwise wrap it in ExceptionInInitializerError.
func WrapInitFailure(loader Loader, thrown *Object) *Object {
	if thrown == nil {
		eiieClass := loader.LoadClass("java/lang/ExceptionInInitializerError")
		return NewObject(eiieClass)
	}
	// Check if thrown is an Error.
	for cls := thrown.Class(); cls != nil; cls = cls.superClass {
		if cls.name == "java/lang/Error" {
			return thrown // propagate directly
		}
	}
	// Wrap in ExceptionInInitializerError.
	eiieClass := loader.LoadClass("java/lang/ExceptionInInitializerError")
	eiie := NewObject(eiieClass)
	setEIIEData(loader, eiie, thrown)
	return eiie
}

// setEIIEData copies the cause's message into the EIIE's detailMessage and
// stores the cause in the EIIE's extra field.
func setEIIEData(loader Loader, eiie, cause *Object) {
	msg := getThrowableMessage(cause)
	for cls := eiie.Class(); cls != nil; cls = cls.superClass {
		if f := cls.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
			if msg != "" {
				strClass := loader.LoadClass("java/lang/String")
				msgObj := NewObject(strClass)
				msgObj.SetExtra(newStringValueFromGo(msg))
				eiie.Cells()[f.SlotID()].SetRef(msgObj)
			}
			break
		}
	}
	eiie.SetExtra(cause)
}

// getThrowableMessage returns the detailMessage string of a Throwable object.
func getThrowableMessage(obj *Object) string {
	for cls := obj.Class(); cls != nil; cls = cls.superClass {
		if f := cls.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
			if msgObj := obj.Cells()[f.SlotID()].GetRef(); msgObj != nil {
				if sv, ok := msgObj.Extra().(*StringValue); ok {
					return sv.GoString()
				}
			}
			return ""
		}
	}
	return ""
}

// setDetailMessage writes a string into a Throwable's detailMessage field.
func setDetailMessage(loader Loader, obj *Object, msg string) {
	for cls := obj.Class(); cls != nil; cls = cls.superClass {
		if f := cls.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
			if msg != "" {
				strClass := loader.LoadClass("java/lang/String")
				strObj := NewObject(strClass)
				strObj.SetExtra(newStringValueFromGo(msg))
				obj.Cells()[f.SlotID()].SetRef(strObj)
			}
			return
		}
	}
}

// newStringValueFromGo converts a Go string to a StringValue by encoding
// each rune as UTF-16 code units. Used for diagnostic messages and class
// names that originate from Go-level strings rather than classfile data.
func newStringValueFromGo(s string) *StringValue {
	if s == "" {
		return NewStringValue([]uint16{})
	}
	ascii := true
	for _, r := range s {
		if r >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		units := make([]uint16, len(s))
		for i, b := range []byte(s) {
			units[i] = uint16(b)
		}
		return NewStringValue(units)
	}
	var units []uint16
	for _, r := range s {
		if r < 0x10000 {
			units = append(units, uint16(r))
		} else {
			r -= 0x10000
			units = append(units, uint16((r>>10)&0x3FF)+0xD800)
			units = append(units, uint16(r&0x3FF)+0xDC00)
		}
	}
	return NewStringValue(units)
}
