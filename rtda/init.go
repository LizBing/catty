// Package rtda — shared class/interface initialization service (ADR-0025,
// ADR-0029).
//
// This file implements the Java 25 cross-context initialization state machine
// per JVMS §5.5, used across Interpreter, IR, and AOT. Each runtime class or
// interface identity has one of four states: not-initialized, initializing,
// initialized, or erroneous.
//
// The per-Class initMu/initCond (ADR-0029) provide the synchronization protocol:
// the lock guards state/owner and the condition is used for other-owner wait
// with notify-all on every terminal transition. Interrupt status of init waiters
// is unchanged (the init wait is not interruptible).
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
// class or interface. It is safe for concurrent use by multiple execution
// contexts — the per-Class initMu/initCond (ADR-0029) implement the
// cross-context synchronization protocol.
//
// Protocol (under initMu):
//   - initialized → return normally (acquire visibility for published state)
//   - same-owner initInProgress → return normally (recursive request)
//   - erroneous → return NoClassDefFoundError
//   - initNotStarted → claim ownership, RELEASE THE LOCK for the <clinit> body,
//     run predecessors + <clinit> via the ClinitRunner callback, then under lock
//     publish the terminal state and Broadcast (notify-all)
//   - other-owner initInProgress → WAIT on initCond WITHOUT holding any Java
//     monitor; on wake re-read state and proceed. Interrupt status of the waiter
//     is UNCHANGED by this wait (ADR-0029).
//
// Terminal publication writes state under initMu and Broadcasts, establishing
// release/acquire visibility for the initialized state and any clinit-written
// static fields (ADR-0030 heap cells).
//
// Parameters:
//   - loader: for resolving predecessor classes
//   - class: the class/interface to initialize
//   - ecID: execution-context identity (for recursive-request recognition)
//   - runClinit: callback that executes <clinit> bytecode
func InitializeClass(loader Loader, class *Class, ecID uint64, runClinit ClinitRunner) InitResult {
	class.initMu.Lock()

	switch class.initState {
	case initInitialized:
		// Already initialized — acquire visibility for published state.
		class.initMu.Unlock()
		return SuccessInit()

	case initInProgress:
		if class.initOwner == ecID {
			// Recursive same-owner request — return normally without
			// re-running <clinit> (JVMS §5.5 step 3).
			class.initMu.Unlock()
			return SuccessInit()
		}
		// Other-owner initInProgress — WAIT on the per-Class condition.
		// The wait is NOT interruptible: interrupt status is unchanged
		// and no InterruptedException is thrown (ADR-0029).
		for class.initState == initInProgress {
			class.initCond.Wait() // releases initMu, reacquires on wake
		}
		// On wake, re-read state and proceed.
		if class.initState == initInitialized {
			class.initMu.Unlock()
			return SuccessInit()
		}
		// Must be erroneous.
		class.initMu.Unlock()
		return buildNCDFE(loader, class.name)

	case initErroneous:
		// Erroneous — throw NoClassDefFoundError (JVMS §5.5 step 2).
		class.initMu.Unlock()
		return buildNCDFE(loader, class.name)

	default: // initNotStarted
		// Claim ownership under the lock.
		class.initState = initInProgress
		class.initOwner = ecID
		class.initMu.Unlock() // RELEASE THE LOCK FOR <clinit> BODY

		// Run predecessor initialization and <clinit> outside the lock
		// (JVMS §5.5 steps 6–7). The owner identity is visible to
		// concurrent callers so recursive same-owner requests return
		// normally.
		var result InitResult
		if !class.IsInterface() {
			if class.superClass != nil {
				if r := InitializeClass(loader, class.superClass, ecID, runClinit); r.ErrObj != nil {
					result = r
				}
			}
			if result.ErrObj == nil {
				for _, iface := range class.DefaultBearingSuperInterfaces(make(map[string]bool)) {
					if r := InitializeClass(loader, iface, ecID, runClinit); r.ErrObj != nil {
						result = r
						break
					}
				}
			}
		}

		if result.ErrObj == nil {
			clinit := class.GetMethod("<clinit>", "()V")
			if clinit != nil {
				result = runClinit(class, clinit)
			}
		}

		// Publish terminal state under lock with notify-all.
		// Terminal publication establishes release/acquire visibility for
		// the initialized state and any clinit-written static fields.
		class.initMu.Lock()
		if result.ErrObj != nil {
			class.initState = initErroneous
		} else {
			class.initState = initInitialized
		}
		class.initOwner = 0
		class.initCond.Broadcast() // wake all waiters (ADR-0029)
		class.initMu.Unlock()

		return result
	}
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
				eiie.SetRefCell(int(f.SlotID()), msgObj)
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
			if msgObj := obj.GetRefCell(int(f.SlotID())); msgObj != nil {
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
				obj.SetRefCell(int(f.SlotID()), strObj)
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
