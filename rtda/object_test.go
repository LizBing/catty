package rtda

import (
	"math"
	"sync"
	"testing"
)

// buildTestObj creates an Object with the given number of heap cells,
// sets typed values in them, and returns the object along with the expected values.
func buildTestObjInts(n int, start int32) *Object {
	cls := &Class{name: "test/Ints"}
	o := &Object{class: cls, heapCells: NewHeapCells(n)}
	for i := 0; i < n; i++ {
		o.SetIntCell(i, start+int32(i))
	}
	return o
}

func buildTestObjLongs(vals ...int64) *Object {
	cls := &Class{name: "test/Longs"}
	o := &Object{class: cls, heapCells: NewHeapCells(len(vals))}
	for i, v := range vals {
		o.SetLongCell(i, v)
	}
	return o
}

func buildTestObjDoubles(vals ...float64) *Object {
	cls := &Class{name: "test/Doubles"}
	o := &Object{class: cls, heapCells: NewHeapCells(len(vals))}
	for i, v := range vals {
		o.SetDoubleCell(i, v)
	}
	return o
}

func buildTestObjRefs(refs ...*Object) *Object {
	cls := &Class{name: "test/Refs"}
	o := &Object{class: cls, heapCells: NewHeapCells(len(refs))}
	for i, r := range refs {
		o.SetRefCell(i, r)
	}
	return o
}

// TestObjectTypedAccessors verifies that Object.GetXxxCell/SetXxxCell
// round-trip all types correctly.
func TestObjectTypedAccessors(t *testing.T) {
	t.Run("int accessors", func(t *testing.T) {
		obj := buildTestObjInts(5, 10)
		for i := 0; i < 5; i++ {
			if got := obj.GetIntCell(i); got != 10+int32(i) {
				t.Errorf("GetIntCell(%d) = %d, want %d", i, got, 10+int32(i))
			}
		}
		obj.SetIntCell(2, 99)
		if got := obj.GetIntCell(2); got != 99 {
			t.Errorf("GetIntCell(2) after set = %d, want 99", got)
		}
	})

	t.Run("long accessors preserve 64 bits", func(t *testing.T) {
		want0 := int64(-6148914691236517200) // 0xAAAAAAAAAAAAAAAA
		obj := buildTestObjLongs(want0, math.MinInt64, math.MaxInt64, 0)
		if got := obj.GetLongCell(0); got != want0 {
			t.Errorf("GetLongCell(0) = 0x%016x, want 0x%016x", uint64(got), uint64(want0))
		}
		if got := obj.GetLongCell(1); got != math.MinInt64 {
			t.Errorf("GetLongCell(1) = %d, want %d", got, math.MinInt64)
		}
		if got := obj.GetLongCell(2); got != math.MaxInt64 {
			t.Errorf("GetLongCell(2) = %d, want %d", got, math.MaxInt64)
		}
		if got := obj.GetLongCell(3); got != 0 {
			t.Errorf("GetLongCell(3) = %d, want 0", got)
		}
	})

	t.Run("float accessors", func(t *testing.T) {
		cls := &Class{name: "test/Floats"}
		obj := &Object{class: cls, heapCells: NewHeapCells(3)}
		obj.SetFloatCell(0, 3.14)
		obj.SetFloatCell(1, float32(math.NaN()))
		obj.SetFloatCell(2, float32(math.Inf(-1)))
		if got := obj.GetFloatCell(0); got != 3.14 {
			t.Errorf("GetFloatCell(0) = %v, want 3.14", got)
		}
		if got := obj.GetFloatCell(1); !math.IsNaN(float64(got)) {
			t.Errorf("GetFloatCell(1) = %v, want NaN", got)
		}
		if got := obj.GetFloatCell(2); !math.IsInf(float64(got), -1) {
			t.Errorf("GetFloatCell(2) = %v, want -Inf", got)
		}
	})

	t.Run("double accessors preserve 64 bits", func(t *testing.T) {
		obj := buildTestObjDoubles(math.Pi, math.E, math.MaxFloat64, math.SmallestNonzeroFloat64)
		if got := obj.GetDoubleCell(0); got != math.Pi {
			t.Errorf("GetDoubleCell(0) = %v, want %v", got, math.Pi)
		}
		if got := obj.GetDoubleCell(1); got != math.E {
			t.Errorf("GetDoubleCell(1) = %v, want %v", got, math.E)
		}
		if got := obj.GetDoubleCell(2); got != math.MaxFloat64 {
			t.Errorf("GetDoubleCell(2) = %v, want %v", got, math.MaxFloat64)
		}
		if got := obj.GetDoubleCell(3); got != math.SmallestNonzeroFloat64 {
			t.Errorf("GetDoubleCell(3) = %v, want %v", got, math.SmallestNonzeroFloat64)
		}
	})

	t.Run("ref accessors", func(t *testing.T) {
		a, b := &Object{}, &Object{}
		obj := buildTestObjRefs(a, b, nil)
		if got := obj.GetRefCell(0); got != a {
			t.Error("GetRefCell(0) != a")
		}
		if got := obj.GetRefCell(1); got != b {
			t.Error("GetRefCell(1) != b")
		}
		if got := obj.GetRefCell(2); got != nil {
			t.Errorf("GetRefCell(2) = %v, want nil", got)
		}
		obj.SetRefCell(2, a)
		if got := obj.GetRefCell(2); got != a {
			t.Error("GetRefCell(2) after set != a")
		}
	})
}

// TestCopyObjectCellsOverlap verifies Java memmove semantics for
// System.arraycopy with overlapping source and destination regions
// within the same array.
func TestCopyObjectCellsOverlap(t *testing.T) {
	t.Run("forward overlap (srcOff < dstOff) ints", func(t *testing.T) {
		// [0, 1, 2, 3, 4] copy src=1..3 to dst=2..4 → [0, 1, 1, 2, 3]
		obj := buildTestObjInts(5, 0) // [0, 1, 2, 3, 4]
		CopyObjectCells(obj, obj, 2, 1, 3, kindInt)
		want := []int32{0, 1, 1, 2, 3}
		for i, w := range want {
			if got := obj.GetIntCell(i); got != w {
				t.Errorf("cell[%d] = %d, want %d", i, got, w)
			}
		}
	})

	t.Run("no overlap (srcOff == dstOff) ints", func(t *testing.T) {
		obj := buildTestObjInts(5, 0) // [0, 1, 2, 3, 4]
		CopyObjectCells(obj, obj, 2, 2, 3, kindInt)
		want := []int32{0, 1, 2, 3, 4}
		for i, w := range want {
			if got := obj.GetIntCell(i); got != w {
				t.Errorf("cell[%d] = %d, want %d (no-op)", i, got, w)
			}
		}
	})

	t.Run("no overlap (srcOff > dstOff, simple forward) ints", func(t *testing.T) {
		// [0, 1, 2, 3, 4] copy src=2..4 to dst=0..2 → [2, 3, 4, 3, 4]
		obj := buildTestObjInts(5, 0)
		CopyObjectCells(obj, obj, 0, 2, 3, kindInt)
		want := []int32{2, 3, 4, 3, 4}
		for i, w := range want {
			if got := obj.GetIntCell(i); got != w {
				t.Errorf("cell[%d] = %d, want %d", i, got, w)
			}
		}
	})

	t.Run("forward overlap longs preserve 64-bit", func(t *testing.T) {
		a := int64(-6148914691236517200) // 0xAAAAAAAAAAAAAAAA
		b := int64(-4919131752989213767) // 0xBBBBBBBBBBBBBBBB
		c := int64(-3689348814741910327) // 0xCCCCCCCCCCCCCCCC
		d := int64(-2459565876494606890) // 0xDDDDDDDDDDDDDDDD
		e := int64(-1229782938247303450) // 0xEEEEEEEEEEEEEEEE
		vals := []int64{a, b, c, d, e}
		obj := buildTestObjLongs(vals...)
		// Copy src=1..4 to dst=2..5 (overlap) — reverse copy required
		CopyObjectCells(obj, obj, 2, 1, 3, kindLong)
		want := []int64{a, b, b, c, d}
		for i, w := range want {
			if got := obj.GetLongCell(i); got != w {
				t.Errorf("cell[%d] = 0x%016x, want 0x%016x", i, uint64(got), uint64(w))
			}
		}
	})

	t.Run("forward overlap doubles preserve 64-bit", func(t *testing.T) {
		vals := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		obj := buildTestObjDoubles(vals...)
		CopyObjectCells(obj, obj, 2, 1, 3, kindDouble)
		want := []float64{1.0, 2.0, 2.0, 3.0, 4.0}
		for i, w := range want {
			if got := obj.GetDoubleCell(i); got != w {
				t.Errorf("cell[%d] = %v, want %v", i, got, w)
			}
		}
	})

	t.Run("cross-object forward copy", func(t *testing.T) {
		src := buildTestObjInts(4, 100)
		dst := buildTestObjInts(6, 0)
		CopyObjectCells(dst, src, 1, 0, 4, kindInt)
		want := []int32{0, 100, 101, 102, 103, 5}
		for i, w := range want {
			if got := dst.GetIntCell(i); got != w {
				t.Errorf("dst[%d] = %d, want %d", i, got, w)
			}
		}
	})

	t.Run("zero length copy is no-op", func(t *testing.T) {
		src := buildTestObjInts(3, 42)
		dst := &Object{class: &Class{name: "test/Dst"}, heapCells: NewHeapCells(3)}
		CopyObjectCells(dst, src, 0, 0, 0, kindInt)
		// all dst cells should still be 0
		for i := 0; i < 3; i++ {
			if got := dst.GetIntCell(i); got != 0 {
				t.Errorf("dst[%d] = %d, want 0 (zero-length copy)", i, got)
			}
		}
	})
}

// TestCloneObject verifies CloneObject does a shallow copy with all field types.
func TestCloneObject(t *testing.T) {
	cls := &Class{name: "test/Mixed"}
	obj := &Object{class: cls, heapCells: NewHeapCells(4)}
	obj.SetIntCell(0, 42)
	obj.SetLongCell(1, math.MaxInt64)
	obj.SetDoubleCell(2, math.Pi)
	inner := &Object{}
	obj.SetRefCell(3, inner)

	clone := CloneObject(obj)
	if clone == obj {
		t.Fatal("CloneObject returned same object, want new instance")
	}
	if clone.class != obj.class {
		t.Error("clone class != original class")
	}
	if len(clone.heapCells) != 4 {
		t.Fatalf("clone has %d cells, want 4", len(clone.heapCells))
	}
	if clone.GetIntCell(0) != 42 {
		t.Errorf("clone int = %d, want 42", clone.GetIntCell(0))
	}
	if clone.GetLongCell(1) != math.MaxInt64 {
		t.Errorf("clone long = %d, want %d", clone.GetLongCell(1), math.MaxInt64)
	}
	if clone.GetDoubleCell(2) != math.Pi {
		t.Errorf("clone double = %v, want %v", clone.GetDoubleCell(2), math.Pi)
	}
	if clone.GetRefCell(3) != inner {
		t.Error("clone ref is shallow: must point to same inner object")
	}
}

// TestConcurrentObjectAccess verifies no data races when multiple goroutines
// read typed cells concurrently.
func TestConcurrentObjectAccess(t *testing.T) {
	t.Run("concurrent typed reads", func(t *testing.T) {
		obj := buildTestObjInts(10, 1)
		var wg sync.WaitGroup
		for g := 0; g < 50; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 10; i++ {
					if obj.GetIntCell(i) != 1+int32(i) {
						t.Errorf("concurrent GetIntCell(%d) = %d, want %d", i, obj.GetIntCell(i), 1+int32(i))
					}
				}
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent typed writes", func(t *testing.T) {
		obj := buildTestObjLongs(0, 0)
		var wg sync.WaitGroup
		for g := 0; g < 50; g++ {
			wg.Add(1)
			go func(id int64) {
				defer wg.Done()
				obj.SetLongCell(0, id)
				obj.SetLongCell(1, -id)
			}(int64(g))
		}
		wg.Wait()
	})
}
