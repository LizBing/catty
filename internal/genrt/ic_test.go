package genrt

import (
	"testing"

	"catty/internal/kernel"
)

// TestEmitterICSlotAgreement pins the emitter's baked inline-cache slot
// literal (emitter.icSlotFor, duplicated to avoid an import cycle) against
// genrt's runtime icHash. A drift would route call sites to the wrong
// cache slot — still CORRECT (slot verification uses full identity) but a
// silent perf cliff.
func TestEmitterICSlotAgreement(t *testing.T) {
	// The emitter's implementation is duplicated; recompute its algorithm
	// here so any change on either side trips this test.
	emitterAlgo := func(cls, name, desc string) uint32 {
		h := uint32(2166136261)
		for _, s := range [3]string{cls, name, desc} {
			for i := 0; i < len(s); i++ {
				h ^= uint32(s[i])
				h *= 16777619
			}
		}
		return h % 1024
	}

	cases := []struct{ cls, name, desc string }{
		{"VCallBench$Node", "get", "()I"},
		{"com/eclipsesource/json/JsonObject", "get", "(Ljava/lang/String;)Lcom/eclipsesource/json/JsonValue;"},
		{"java/util/HashMap", "put", "(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"},
		{"", "", ""},
		{"a", "b", "c"},
	}
	for _, c := range cases {
		if got := ICSlot(c.cls, c.name, c.desc); got != emitterAlgo(c.cls, c.name, c.desc) {
			t.Errorf("ICSlot(%q,%q,%q) = %d, emitter algo = %d",
				c.cls, c.name, c.desc, got, emitterAlgo(c.cls, c.name, c.desc))
		}
		if ICSlot(c.cls, c.name, c.desc) >= 1024 {
			t.Fatalf("slot out of range: %d", ICSlot(c.cls, c.name, c.desc))
		}
	}
}

// TestInlineCacheHitAndBench verifies the monomorphic cache returns the
// same method as the uncached walk and survives kernel reinstall.
func TestInlineCacheHitAndReinstall(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: discardWriter{}, Stderr: discardWriter{}})
	InstallKernel(k)
	th := &testOwner{key: k.MintKey()}

	obj := New(th, "java/util/HashMap")
	CallSpecial(th, obj, "java/util/HashMap", "<init>", "()V", []kernel.Value{})

	slot := ICSlot("java/util/HashMap", "size", "()I")
	m1 := methodForDynCachedSlot(slot, ClassOf(obj), "java/util/HashMap", "size", "()I")
	m2 := methodForDynCachedSlot(slot, ClassOf(obj), "java/util/HashMap", "size", "()I")
	if m1 != m2 {
		t.Fatal("monomorphic cache missed on repeat dispatch")
	}
	e := icTable[slot].Load()
	if e == nil || e.dyn != ClassOf(obj) {
		t.Fatal("cache entry not populated")
	}

	// Reinstall must invalidate (stale kernels must never leak through).
	k2 := kernel.New(kernel.Options{Stdout: discardWriter{}, Stderr: discardWriter{}})
	InstallKernel(k2)
	if icTable[slot].Load() != nil {
		t.Fatal("icTable not invalidated by InstallKernel")
	}
}
