package kernel

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// Header is embedded by every heap object (instances, arrays, strings).
// Go GC traces all of it; identity semantics are pointer equality.
type Header struct {
	Class *Class
	mu    *sync.Mutex // lazy monitor (synchronized support, M0 stub)
	hash  uint32      // lazy identity hash; 0 = unassigned
}

// Monitor returns the object's monitor, allocating on first use.
// M0: plain mutex under single-threaded execution; wait/notify and
// thin-lock semantics arrive with the M1 monitor rework (P-0001 scope note).
func (h *Header) Monitor() *sync.Mutex {
	if h.mu == nil {
		h.mu = new(sync.Mutex)
	}
	return h.mu
}

var hashCtr atomic.Uint64

// IdentityHash implements java.lang.Object.hashCode() (identity hash).
func (h *Header) IdentityHash() uint32 {
	for {
		v := atomic.LoadUint32(&h.hash)
		if v != 0 {
			return v
		}
		n := uint32(hashCtr.Add(1)) // sequential, stable within a run
		if atomic.CompareAndSwapUint32(&h.hash, 0, n) {
			return n
		}
	}
}

// Instance is an object of a synthesized or interpreted class.
type Instance struct {
	Header
	Fields  []Value // absolute slot layout across the hierarchy (super first)
	Payload any     // scratch storage for synthesized natives (ArrayList, StringBuilder, …)
}

// fieldByName finds an instance field slot value by name across the
// hierarchy (first match wins). Returns nil when absent.
func (in *Instance) fieldByName(name string) Value {
	for x := in.Class; x != nil; x = x.Super {
		for _, f := range x.OwnFields {
			if f.Name == name {
				return in.Fields[f.Slot]
			}
		}
	}
	return nil
}

// ArrayObj is a Java array. Elems are uniform Values in M0 (specialized
// primitive backing stores are a later optimization).
type ArrayObj struct {
	Header
	CompDesc string // "I", "J", "Ljava/lang/Object;", "[I", …
	Elems    []Value
}

// JString is a java/lang.String: UTF-16 code units (JDK 8 semantics,
// no compact strings). Always interned-eligible via Kernel.
type JString struct {
	Header
	Chars []uint16
}

// String converts to a Go string (allocation; natives use sparingly).
func (s *JString) String() string {
	return string(utf16Decode(s.Chars))
}

// AsJString extracts the chars of a Value that must be a String.
func AsJString(v Value) (*JString, error) {
	s, ok := v.(*JString)
	if !ok {
		return nil, fmt.Errorf("value is %T, want *JString", v)
	}
	return s, nil
}

func utf16Decode(u []uint16) []rune {
	r := make([]rune, 0, len(u))
	for i := 0; i < len(u); i++ {
		c := u[i]
		switch {
		case c >= 0xD800 && c <= 0xDBFF && i+1 < len(u) &&
			u[i+1] >= 0xDC00 && u[i+1] <= 0xDFFF:
			r = append(r, ((rune(c)-0xD800)<<10|(rune(u[i+1])-0xDC00))+0x10000)
			i++
		default:
			r = append(r, rune(c))
		}
	}
	return r
}

func utf16Encode(rs []rune) []uint16 {
	out := make([]uint16, 0, len(rs))
	for _, r := range rs {
		if r >= 0x10000 {
			r -= 0x10000
			out = append(out, uint16(0xD800+(r>>10)), uint16(0xDC00+(r&0x3FF)))
		} else {
			out = append(out, uint16(r))
		}
	}
	return out
}

// zeroValue returns the default value for an array component descriptor.
func zeroValue(compDesc string) Value {
	switch compDesc {
	case "J":
		return int64(0)
	case "F":
		return float32(0)
	case "D":
		return float64(0)
	}
	if strings.HasPrefix(compDesc, "L") || strings.HasPrefix(compDesc, "[") {
		return nil // reference component
	}
	return int32(0) // B C S I Z
}
