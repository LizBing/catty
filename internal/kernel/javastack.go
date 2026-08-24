// Java-level stack backfill (DEBT-0019 diagnostic infrastructure).
//
// Both engines route every Java-to-Java call through Kernel.InvokeAs, so
// that single point maintains a per-thread Java call stack. Throwable
// creation sites (bootstrap helpers, Throwable.<init> natives) snapshot it
// onto the throwable, mirroring HotSpot's fillInStackTrace semantics: the
// trace reflects where the exception was CONSTRUCTED, not where it landed.
package kernel

// JavaFrame is one resolved Java call frame. Class is an internal name
// ("com/eclipsesource/json/JsonParser"); Method excludes the descriptor.
type JavaFrame struct {
	Class  string
	Method string
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
//		at com/foo/Bar.baz(Unknown Source)
//
// Frame lines are omitted when no trace was captured (owner-less paths).
// Line numbers are v1 "(Unknown Source)": per-frame call-site pcs are not
// tracked yet (emitter-abi.md §stack-backfill).
func FormatUncaught(threadName string, th *Thrown) string {
	header := "Exception in thread \"" + threadName + "\" "
	msg := ""
	if s, ok := th.Obj.fieldByName("detailMessage").(*JString); ok && s != nil {
		msg = ": " + s.String()
	}
	out := header + dotted(th.Obj.Class.Name) + msg + "\n"
	if th.Obj.Stack != nil {
		// JVM order is leaf-first; our stack stores top-of-stack last.
		fr := th.Obj.Stack.Frames
		for i := len(fr) - 1; i >= 0; i-- {
			out += "\tat " + dotted(fr[i].Class) + "." + fr[i].Method + "(Unknown Source)\n"
		}
	}
	return out
}
