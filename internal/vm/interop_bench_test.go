package vm

import (
	"io"
	"testing"

	"catty/internal/kernel"
)

// BenchmarkInteropNativeCall measures the full dispatch cost of one
// emitted-code → Go-native call (R-0008 spike B): CallVirtualIC →
// invokeChecked (FrameMeter) → InvokeAs (frame push/pop, ctx pool) →
// NativeFunc wrapper → the embedder's Go function. This is the "interop
// tax" an embedding pays per boundary crossing, comparable to the ~70ns
// virtual-call figure in R-0006.
func BenchmarkInteropNativeCall(b *testing.B) {
	k := kernel.New(kernel.Options{Stdout: io.Discard})
	registerTestBridge(k)

	cls, ok := k.ClassByName("go/Bridge")
	if !ok {
		b.Fatal("go/Bridge not registered")
	}
	m, err := k.ResolveMethod(cls, "add", "(II)I")
	if err != nil {
		b.Fatal(err)
	}
	th := New(k)
	a := int32(20)
	c := int32(22)
	args := []kernel.Value{a, c}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := th.Call(m, nil, args); err != nil {
			b.Fatal(err)
		}
	}
}
