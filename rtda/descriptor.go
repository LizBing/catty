package rtda

import "strings"

// MethodDescriptor is the parsed form of a method descriptor (JVMS §4.3.3),
// e.g. "(II)V" -> {parameterTypes:["I","I"], returnType:"V"}.
type MethodDescriptor struct {
	ParameterTypes []string
	ReturnType     string
}

func ParseMethodDescriptor(descriptor string) MethodDescriptor {
	md := MethodDescriptor{}
	start := strings.IndexByte(descriptor, '(')
	end := strings.LastIndexByte(descriptor, ')')
	if start < 0 || end < 0 {
		return md
	}
	raw := descriptor[start+1 : end]
	md.ReturnType = descriptor[end+1:]
	md.ParameterTypes = parseTypeList(raw)
	return md
}

// parseTypeList splits a concatenation of type descriptors into individual ones.
// A type descriptor is one of: a base type single char (B C D F I J S Z V), an
// object type "L...;", or an array type "["+component.
func parseTypeList(raw string) []string {
	var types []string
	for len(raw) > 0 {
		t, rest := parseNextType(raw)
		types = append(types, t)
		raw = rest
	}
	return types
}

func parseNextType(raw string) (typ string, rest string) {
	switch raw[0] {
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 'V':
		return raw[:1], raw[1:]
	case 'L':
		semi := strings.IndexByte(raw, ';')
		return raw[:semi+1], raw[semi+1:]
	case '[':
		// array: consume all leading '[' then the component type
		i := 1
		for i < len(raw) && raw[i] == '[' {
			i++
		}
		comp, _ := parseNextType(raw[i:])
		return raw[:i+len(comp)], raw[i+len(comp):]
	default:
		return raw[:1], raw[1:]
	}
}

// ArgSlots returns the number of local-variable slots the parameters of a method
// occupy: long/double (J, D) take 2 slots, every other type takes 1. Used by the
// classloader (frame sizing) and the lowering pass (invoke slot effects).
func (md MethodDescriptor) ArgSlots() int {
	n := 0
	for _, t := range md.ParameterTypes {
		if t == "J" || t == "D" {
			n += 2
		} else {
			n++
		}
	}
	return n
}
