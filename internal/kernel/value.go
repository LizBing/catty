// Package kernel implements Catty's runtime object model: classes,
// objects, arrays, strings, monitors and the bootstrap (synthesized)
// subset of java.* used by the M0 interpreter.
//
// Layering rule (architecture-rules R7): kernel must not import the VM
// or any higher layer. Callbacks upward go through the Invoker interface.
package kernel

import "strconv"

// Value is the universal interpreter slot.
//
// Conventions:
//   - nil                  → Java null reference
//   - int32                → boolean, byte, char, short, int (per JVM spec,
//     these are all computed as ints; char is zero-extended on load)
//   - int64, float32, float64 → category-2 and FP types
//   - *Instance            → instances of synthesized/interpreted classes
//   - *ArrayObj            → arrays
//   - *JString             → java/lang/String
type Value = any

// TopSentinel fills the second raw slot of category-2 values so that
// stack manipulation opcodes (dup2/pop2/…) can follow JVMS semantics
// verbatim with one Value per slot.
var TopSentinel = new(struct{ reserved byte })

// IsTop reports whether v is a category-2 filler slot.
func IsTop(v Value) bool { return v == Value(TopSentinel) }

// IsCat2 reports whether v occupies two stack/local slots.
func IsCat2(v Value) bool {
	switch v.(type) {
	case int64, float64:
		return true
	}
	return false
}

// ParseMethodDesc splits a method descriptor (JVMS §4.3.3) into argument
// type descriptors and return type descriptor. Example:
//
//	"(IJ[Ljava/lang/String;)V" → ["I","J","[Ljava/lang/String;"], "V"
func ParseMethodDesc(desc string) (args []string, ret string, err error) {
	if len(desc) == 0 || desc[0] != '(' {
		return nil, "", &DescError{Desc: desc, Why: "missing '('"}
	}
	i := 1
	for i < len(desc) && desc[i] != ')' {
		start := i
		for i < len(desc) && desc[i] == '[' {
			i++
		}
		if i >= len(desc) {
			return nil, "", &DescError{Desc: desc, Why: "unterminated array type"}
		}
		switch desc[i] {
		case 'L':
			end := indexSemicolon(desc, i+1)
			if end < 0 {
				return nil, "", &DescError{Desc: desc, Why: "unterminated L type"}
			}
			i = end + 1
		case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z':
			i++
		default:
			return nil, "", &DescError{Desc: desc, Why: "bad type char '" + string(desc[i]) + "'"}
		}
		args = append(args, desc[start:i])
	}
	if i >= len(desc) || desc[i] != ')' {
		return nil, "", &DescError{Desc: desc, Why: "missing ')'"}
	}
	ret = desc[i+1:]
	if ret == "" {
		return nil, "", &DescError{Desc: desc, Why: "missing return type"}
	}
	return args, ret, nil
}

func indexSemicolon(s string, from int) int {
	for j := from; j < len(s); j++ {
		if s[j] == ';' {
			return j
		}
	}
	return -1
}

// SlotCount returns the operand-stack slots a descriptor of an argument
// type consumes (long/double = 2, everything else = 1).
func SlotCount(typeDesc string) int {
	if typeDesc == "J" || typeDesc == "D" {
		return 2
	}
	return 1
}

// DescError reports a malformed type/method descriptor.
type DescError struct {
	Desc string
	Why  string
}

func (e *DescError) Error() string { return "bad descriptor " + strconv.Quote(e.Desc) + ": " + e.Why }

// FormatInt formats per Java Integer.toString semantics.
func FormatInt(v int32) string { return strconv.FormatInt(int64(v), 10) }

// FormatLong formats per Java Long.toString semantics.
func FormatLong(v int64) string { return strconv.FormatInt(v, 10) }
