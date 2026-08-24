// Java-level stack backfill (DEBT-0019 diagnostic infrastructure).
//
// Both engines route every Java-to-Java call through Kernel.InvokeAs, so
// that single point maintains a per-thread Java call stack. Throwable
// creation sites (bootstrap helpers, Throwable.<init> natives) snapshot it
// onto the throwable, mirroring HotSpot's fillInStackTrace semantics: the
// trace reflects where the exception was CONSTRUCTED, not where it landed.
package kernel

import (
	"fmt"
)

// JavaFrame is one resolved Java call frame. Class is an internal name
// ("com/eclipsesource/json/JsonParser"); Method excludes the descriptor.
//
// Line carries the source line the frame was last known to execute:
// for the leaf it is the throw/creation site; for callers it is the line
// of the call into the next frame — exactly what a JVM trace prints.
// Zero means unknown ("Unknown Source").
type JavaFrame struct {
	Class  string
	Method string
	Line   int32
}

// JavaStackTrace is the creation-time Java call stack captured on a
// throwable. Stored behind a pointer so non-throwable instances pay only
// one word.
type JavaStackTrace struct {
	Frames []JavaFrame // outermost first; top of stack last
}

// FrameTracker is implemented by thread identities that carry the Java
// call stack. The VM's Thread implements it; owner-less probes simply
// never satisfy the assertion and skip frame bookkeeping.
type FrameTracker interface {
	PushJavaFrame(JavaFrame)
	PopJavaFrame()
	// JavaFrames returns a snapshot copy, outermost first.
	JavaFrames() []JavaFrame
	// SetTopJavaLine records the source line the top frame is now
	// executing (call sites and line-segment entries; no-op when empty).
	SetTopJavaLine(line int32)
}

// attachTrace snapshots the current Java frames of owner into obj,
// trimming the trailing <init> construction chain exactly like
// fillInStackTrace does (the exception's own constructors do not appear).
func attachTrace(owner OwnerKey, obj *Instance) {
	ft, ok := owner.(FrameTracker)
	if !ok {
		return
	}
	fr := ft.JavaFrames()
	n := len(fr)
	for n > 0 && fr[n-1].Method == "<init>" {
		n--
	}
	cp := make([]JavaFrame, n)
	copy(cp, fr[:n])
	obj.Stack = &JavaStackTrace{Frames: cp}
}

// AttachTraceTo is the exported capture point for engines that create
// bootstrap throwables outside Throwable.<init> (interpreter helpers and
// the genrt bridge pass their thread identity here).
func AttachTraceTo(owner OwnerKey, obj *Instance) { attachTrace(owner, obj) }

// FormatUncaught renders the standard uncaught-throwable report:
//
//	Exception in thread "main" java.lang.NullPointerException: msg
//		at com.foo.Bar.baz(Foo.java:12)
//
// Frame lines render JVM-style "(Source.java:N)" when the frame's line is
// known AND the declaring class carries a SourceFile attribute; otherwise
// "(Unknown Source)". Leaf-first ordering matches the JVM.
func (k *Kernel) FormatUncaught(threadName string, th *Thrown) string {
	header := "Exception in thread \"" + threadName + "\" "
	msg := ""
	if s, ok := th.Obj.fieldByName("detailMessage").(*JString); ok && s != nil {
		msg = ": " + s.String()
	}
	out := header + dotted(th.Obj.Class.Name) + msg + "\n"
	if th.Obj.Stack == nil {
		return out
	}
	fr := th.Obj.Stack.Frames
	for i := len(fr) - 1; i >= 0; i-- {
		f := fr[i]
		at := "\tat " + dotted(f.Class) + "." + f.Method
		if f.Line > 0 {
			if sf := k.sourceFileOf(f.Class); sf != "" {
				at += "(" + sf + ":" + fmt.Sprintf("%d", f.Line) + ")"
			} else {
				at += "(Unknown Source)"
			}
		} else {
			at += "(Unknown Source)"
		}
		out += at + "\n"
	}
	return out
}

func (k *Kernel) sourceFileOf(internalClass string) string {
	c, ok := k.ClassByName(internalClass)
	if !ok || c.CF == nil {
		return ""
	}
	return c.CF.SourceFile
}
