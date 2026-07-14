package interpreter

import "catty/rtda"

// RunMethod runs an INTERPRETED method to completion and returns its result —
// the entry point the AOT bridge (catty/runtime) calls for interpreted targets
// (constructors, user methods). Native targets are handled by the runtime's
// runNative; this is for bytecode.
//
// It sets a bridge-return slot so the method's return (which has no caller frame
// in this context) is captured instead of dropped, then runs the dispatch loop.
func RunMethod(thread *rtda.Thread, method *rtda.Method, args []rtda.Slot) rtda.Slot {
	var ret rtda.Slot
	thread.SetBridgeReturn(&ret)
	frame := thread.NewFrame(method)
	for i, a := range args {
		frame.SetSlot(i, a)
	}
	frame.EnterSyncMonitor()
	thread.PushFrame(frame)
	Loop(thread)
	thread.SetBridgeReturn(nil)
	return ret
}
