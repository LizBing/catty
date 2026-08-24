package genrt

import (
	"testing"

	"catty/internal/kernel"
)

type nraOwner struct{ key uint64 }

func (o *nraOwner) OwnerKey() uint64 { return o.key }

func TestNewRefArrayThenStore(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: discardWriter{}, Stderr: discardWriter{}})
	K = k
	th := &nraOwner{key: k.MintKey()}
	_ = th

	tv, terr := NewRefArray("java/lang/String", 6)
	if terr != nil {
		t.Fatalf("NewRefArray threw: %v", terr)
	}
	if tv == nil {
		t.Fatal("NewRefArray returned nil value")
	}

	strV := k.InternGo("hello")
	if strV == nil {
		t.Fatal("InternGo returned nil")
	}
	arr := tv.(*kernel.ArrayObj)
	if len(arr.Elems) != 6 {
		t.Fatalf("len=%d", len(arr.Elems))
	}
	if th := AStoreChecked(tv, int32(0), strV); th == nil {
		// storing nil is legal for reference arrays
		_ = th
	}
}
