package transpile

import "strings"

// mangle turns a (className, methodName) pair into a valid, collision-resistant
// Go identifier. Internal class names use '/'; both '/' and the angle brackets
// of <init>/<clinit> are replaced with '_'. The class prefix makes collisions
// across classes unlikely and keeps a bare method name from clashing with Go's
// reserved `init`.
func mangle(class, method string) string {
	var b strings.Builder
	for _, r := range class + "_" + method {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
