package emitter

import (
	"fmt"
	"strings"
)

// MangleClass maps an internal class name to a Go identifier fragment
// (emitter-abi §1): '/'→'_', '$'→'_00036'.
func MangleClass(internal string) string {
	s := strings.ReplaceAll(internal, "/", "_")
	return strings.ReplaceAll(s, "$", "_00036")
}

// MangleDesc maps a method descriptor to its identifier fragment:
// parentheses dropped, 'L…;'→'L…_', '['→'A', primitive letters kept.
func MangleDesc(desc string) string {
	var b strings.Builder
	for i := 0; i < len(desc); i++ {
		switch desc[i] {
		case '(', ')', ';', '/':
			b.WriteByte('_')
		case '[':
			b.WriteByte('A')
		default:
			b.WriteByte(desc[i])
		}
	}
	return b.String()
}

// MethodFunc renders the Go function identifier for a method.
func MethodFunc(class, name, desc string) string {
	return fmt.Sprintf("Catty_%s_%s__%s", MangleClass(class), sanitizeName(name), MangleDesc(desc))
}

// sanitizeName keeps <init>/<clinit> out of identifiers.
func sanitizeName(name string) string {
	switch name {
	case "<init>":
		return "init"
	case "<clinit>":
		return "clinit"
	}
	return name
}
