package verify

import (
	"fmt"
	"os"
	"strings"
)

// Verification value types (JVMS §4.10.1 verification_type_info lifted to a
// simulation model). Category-2 values occupy two raw slots: [value][top].

type vkind = byte

const (
	kTop byte = iota
	kInt      // computational type one: boolean/byte/char/short/int
	kFloat
	kLong
	kDouble
	kNull
	kUninitThis
	kObj
	kUninit
	kUnknown // wildcard: produced where knowledge is deferred (DEV-0001)
)

type vtype struct {
	kind  vkind
	class string // kObj: internal class name or array descriptor "[…"
	off   uint16 // kUninit: pc of the `new` that created it
}

var (
	tTop     = vtype{kind: kTop}
	tInt     = vtype{kind: kInt}
	tFloat   = vtype{kind: kFloat}
	tLong    = vtype{kind: kLong}
	tDouble  = vtype{kind: kDouble}
	tNull    = vtype{kind: kNull}
	tUninitT = vtype{kind: kUninitThis}
	tUnknown = vtype{kind: kUnknown}
)

func tObj(class string) vtype { return vtype{kind: kObj, class: class} }
func tUninit(off uint16) vtype {
	return vtype{kind: kUninit, off: off}
}

func isCat2Kind(k vkind) bool { return k == kLong || k == kDouble }
func isCat2(v vtype) bool     { return isCat2Kind(v.kind) }

func isRef(v vtype) bool {
	switch v.kind {
	case kNull, kUninitThis, kObj, kUninit, kUnknown:
		return true
	}
	return false
}

// String renders a type for diagnostics.
func (v vtype) String() string {
	switch v.kind {
	case kTop:
		return "top"
	case kInt:
		return "int"
	case kFloat:
		return "float"
	case kLong:
		return "long"
	case kDouble:
		return "double"
	case kNull:
		return "null"
	case kUninitThis:
		return "uninitThis"
	case kObj:
		return "object:" + v.class
	case kUninit:
		return fmt.Sprintf("uninit@%d", v.off)
	case kUnknown:
		return "unknown"
	}
	return "?"
}

// refCompatible decides child <: parent for reference kinds, consulting the
// resolver over the loaded class graph; unknown names defer (accept).
func (c *checker) refCompatible(child, parent string) bool {
	if child == parent {
		return true
	}
	if parent == "java/lang/Object" || parent == "java/io/Serializable" {
		// Object is the universal supertype; Serializable is a marker
		// interface every value type satisfies (javac frames type it
		// liberally — e.g. String/Integer merges in minimal-json).
		return true
	}
	if strings.HasPrefix(child, "[") || strings.HasPrefix(parent, "[") {
		return c.arrayCompatible(child, parent)
	}
	if !c.r.Known(child) || !c.r.Known(parent) {
		return true // DEV-0001 deferred strictness
	}
	return c.r.IsSubclass(child, parent)
}

func (c *checker) arrayCompatible(childDesc, parentDesc string) bool {
	if childDesc == parentDesc {
		return true
	}
	if parentDesc == "java/lang/Object" {
		return true
	}
	if !strings.HasPrefix(parentDesc, "[") {
		return false // primitive arrays are Object-only
	}
	childComp := childDesc[1:]
	parentComp := parentDesc[1:]
	// X[] <: Y[] iff X <: Y for reference components (covariance)
	if strings.HasPrefix(childComp, "L") && strings.HasPrefix(parentComp, "L") {
		return c.refCompatible(childComp[1:len(childComp)-1], parentComp[1:len(parentComp)-1])
	}
	if strings.HasPrefix(childComp, "[") && strings.HasPrefix(parentComp, "[") {
		return c.arrayCompatible(childComp, parentComp)
	}
	if !strings.HasPrefix(childComp, "L") && !strings.HasPrefix(childComp, "[") {
		return false // primitive component arrays are invariant vs ref arrays
	}
	// child is reference-component array, parent component is unknown-shape
	return c.refCompatible(childComp, parentComp)
}

// compatible checks got ⊑ want at merge points. Unknowns accept; nulls
// accept into any reference slot.
func (c *checker) itemCompatible(got, want vtype) error {
	if want.kind == kTop {
		if got.kind != kTop {
			return fmt.Errorf("want top, got %s", got)
		}
		return nil
	}
	if got.kind == kUnknown || want.kind == kUnknown {
		return nil
	}
	switch want.kind {
	case kInt:
		if got.kind != kInt {
			return fmt.Errorf("want int, got %s", got)
		}
	case kFloat:
		if got.kind != kFloat {
			return fmt.Errorf("want float, got %s", got)
		}
	case kLong:
		if got.kind != kLong {
			return fmt.Errorf("want long, got %s", got)
		}
	case kDouble:
		if got.kind != kDouble {
			return fmt.Errorf("want double, got %s", got)
		}
	case kNull:
		if got.kind != kNull {
			return fmt.Errorf("want null, got %s", got)
		}
	case kUninitThis:
		if got.kind != kUninitThis && got.kind != kNull && got.kind != kObj {
			return fmt.Errorf("want uninitThis-compatible, got %s", got)
		}
	case kObj:
		if !isRef(got) {
			return fmt.Errorf("want reference, got %s", got)
		}
		switch got.kind {
		case kNull, kUninitThis:
			return nil
		case kUninit:
			return nil // javac types pre-init receivers as the holder class
		}
		gotClass := got.class
		if !c.refCompatible(gotClass, want.class) {
			return fmt.Errorf("%s not assignable to %s", gotClass, want.class)
		}
	case kUninit:
		if got.kind != kUninit || got.off != want.off {
			return fmt.Errorf("want uninit@%d, got %s", want.off, got)
		}
	}
	return nil
}

// stateCompatible checks a full simulated state against an SM frame.
func (c *checker) stateCompatible(st *frame, fr canonicalFrame) error {
	// Locals: the frame constrains only its declared prefix; deeper
	// simulated locals are dead values and unconstrained (JVMS §4.10.1).
	if len(st.locals) < len(fr.locals) {
		return fmt.Errorf("locals depth %d, frame needs %d", len(st.locals), len(fr.locals))
	}
	for i := range fr.locals {
		if err := c.itemCompatible(st.locals[i], fr.locals[i]); err != nil {
			if os.Getenv("CATTY_VTRACE") != "" {
				fmt.Fprintf(os.Stderr, "[vtrace] %s local[%d]: got=%v want=%v\n", c.mn, i, st.locals[i], fr.locals[i])
			}
			return fmt.Errorf("local[%d]: %w", i, err)
		}
	}
	if len(st.stack) != len(fr.stack) {
		return fmt.Errorf("stack depth %d, frame says %d", len(st.stack), len(fr.stack))
	}
	for i := range fr.stack {
		if err := c.itemCompatible(st.stack[i], fr.stack[i]); err != nil {
			return fmt.Errorf("stack[%d]: %w", i, err)
		}
	}
	return nil
}

// arrayElem returns the element type of an array descriptor.
func arrayElem(desc string) vtype {
	comp := desc[1:]
	switch comp {
	case "B", "C", "S", "I", "Z":
		return tInt
	case "J":
		return tLong
	case "F":
		return tFloat
	case "D":
		return tDouble
	}
	if comp[0] == 'L' {
		return tObj(comp[1 : len(comp)-1])
	}
	return tObj(comp)
}
