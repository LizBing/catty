package rtda

import (
	"math"
	"sync"
	"testing"
)

// TestHeapCellTypedAccess verifies that each typed getter/setter round-trips
// the full bit pattern, including 64-bit long and double values.
func TestHeapCellTypedAccess(t *testing.T) {
	t.Run("int32 round-trip", func(t *testing.T) {
		c := &HeapCell{}
		c.SetInt(-42)
		if got := c.GetInt(); got != -42 {
			t.Errorf("GetInt() = %d, want -42", got)
		}
	})

	t.Run("float32 round-trip preserves NaN", func(t *testing.T) {
		c := &HeapCell{}
		v := float32(math.NaN())
		c.SetFloat(v)
		got := c.GetFloat()
		if !math.IsNaN(float64(got)) {
			t.Errorf("GetFloat() = %v, want NaN", got)
		}
		// signalled NaN bit pattern
		v2 := float32(math.Float32frombits(0x7fc00000))
		c.SetFloat(v2)
		if got2 := c.GetFloat(); math.Float32bits(got2) != 0x7fc00000 {
			t.Errorf("GetFloat() bits = 0x%08x, want 0x7fc00000", math.Float32bits(got2))
		}
	})

	t.Run("int64 long round-trip", func(t *testing.T) {
		c := &HeapCell{}
		// Value that uses both halves — catches 32-bit truncation
		v := int64(-6148914691236517200) // bits: 0xAAAAAAAA55555555
		c.SetLong(v)
		if got := c.GetLong(); got != v {
			t.Errorf("GetLong() = 0x%016x, want 0x%016x", uint64(got), uint64(v))
		}
	})

	t.Run("int64 long min/max", func(t *testing.T) {
		c := &HeapCell{}
		for _, v := range []int64{0, -1, math.MinInt64, math.MaxInt64, 1, -2} {
			c.SetLong(v)
			if got := c.GetLong(); got != v {
				t.Errorf("GetLong() = %d, want %d", got, v)
			}
		}
	})

	t.Run("float64 double round-trip", func(t *testing.T) {
		c := &HeapCell{}
		for _, v := range []float64{0.0, -0.0, 1.0, math.Pi, math.MaxFloat64,
			math.SmallestNonzeroFloat64, math.Inf(1), math.Inf(-1)} {
			c.SetDouble(v)
			got := c.GetDouble()
			if math.IsNaN(v) {
				if !math.IsNaN(got) {
					t.Errorf("GetDouble() = %v, want NaN", got)
				}
			} else if got != v {
				t.Errorf("GetDouble() = %v, want %v", got, v)
			}
		}
	})

	t.Run("float64 double NaN preserves bits", func(t *testing.T) {
		c := &HeapCell{}
		bits := uint64(0x7ff8000000000001)
		v := math.Float64frombits(bits)
		c.SetDouble(v)
		got := c.GetDouble()
		if math.Float64bits(got) != bits {
			t.Errorf("GetDouble() bits = 0x%016x, want 0x%016x", math.Float64bits(got), bits)
		}
	})

	t.Run("reference round-trip", func(t *testing.T) {
		c := &HeapCell{}
		if got := c.GetRef(); got != nil {
			t.Errorf("GetRef() on zero cell = %v, want nil", got)
		}
		obj := &Object{}
		c.SetRef(obj)
		if got := c.GetRef(); got != obj {
			t.Errorf("GetRef() = %v, want %v", got, obj)
		}
		c.SetRef(nil)
		if got := c.GetRef(); got != nil {
			t.Errorf("GetRef() after SetRef(nil) = %v, want nil", got)
		}
	})

	t.Run("type isolation", func(t *testing.T) {
		// Writing a long should not corrupt a subsequent float read,
		// and vice versa.
		c := &HeapCell{}
		c.SetLong(int64(-6148914691236517200)) // high bits set
		c.SetFloat(3.14)
		c.SetLong(int64(-8198552921648689601)) // different high bits
		want := int64(-8198552921648689601)
		if got := c.GetLong(); got != want {
			t.Errorf("long after float overwrite = 0x%016x, want 0x%016x", uint64(got), uint64(want))
		}
	})
}

// TestHeapCellConcurrentAccess verifies that concurrent reads/writes of heap
// cells do not produce data races (verified by the -race flag).
func TestHeapCellConcurrentAccess(t *testing.T) {
	t.Run("concurrent int writes", func(t *testing.T) {
		c := &HeapCell{}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(v int32) {
				c.SetInt(v)
				_ = c.GetInt()
				wg.Done()
			}(int32(i))
		}
		wg.Wait()
	})

	t.Run("concurrent ref writes", func(t *testing.T) {
		c := &HeapCell{}
		obj := &Object{}
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				c.SetRef(obj)
				_ = c.GetRef()
				c.SetRef(nil)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

// TestNewHeapCells validates bulk allocation.
func TestNewHeapCells(t *testing.T) {
	cells := NewHeapCells(5)
	if len(cells) != 5 {
		t.Fatalf("NewHeapCells(5) len = %d, want 5", len(cells))
	}
	// All cells should be zero-initialized
	for i := range cells {
		if cells[i].GetInt() != 0 {
			t.Errorf("cells[%d].GetInt() = %d, want 0", i, cells[i].GetInt())
		}
		if cells[i].GetRef() != nil {
			t.Errorf("cells[%d].GetRef() = %v, want nil", i, cells[i].GetRef())
		}
	}
}

// TestCopyHeapCells verifies typed copy with 64-bit preservation.
func TestCopyHeapCells(t *testing.T) {
	t.Run("copy longs", func(t *testing.T) {
		src := NewHeapCells(3)
		w0 := int64(-6148914691236517200)
		w1 := int64(-4919131752989213767)
		w2 := int64(-3689348814741910327)
		src[0].SetLong(w0)
		src[1].SetLong(w1)
		src[2].SetLong(w2)
		dst := NewHeapCells(3)
		CopyHeapCells(dst, src)
		for i, want := range []int64{w0, w1, w2} {
			if got := dst[i].GetLong(); got != want {
				t.Errorf("dst[%d].GetLong() = 0x%016x, want 0x%016x", i, uint64(got), uint64(want))
			}
		}
	})

	t.Run("copy doubles", func(t *testing.T) {
		src := NewHeapCells(2)
		src[0].SetDouble(math.Pi)
		src[1].SetDouble(math.E)
		dst := NewHeapCells(2)
		CopyHeapCells(dst, src)
		if got := dst[0].GetDouble(); got != math.Pi {
			t.Errorf("dst[0].GetDouble() = %v, want %v", got, math.Pi)
		}
		if got := dst[1].GetDouble(); got != math.E {
			t.Errorf("dst[1].GetDouble() = %v, want %v", got, math.E)
		}
	})

	t.Run("copy refs", func(t *testing.T) {
		src := NewHeapCells(2)
		obj1, obj2 := &Object{}, &Object{}
		src[0].SetRef(obj1)
		src[1].SetRef(obj2)
		dst := NewHeapCells(2)
		CopyHeapCells(dst, src)
		if dst[0].GetRef() != obj1 {
			t.Error("dst[0] ref not preserved")
		}
		if dst[1].GetRef() != obj2 {
			t.Error("dst[1] ref not preserved")
		}
	})

	t.Run("dst shorter than src", func(t *testing.T) {
		src := NewHeapCells(5)
		src[0].SetInt(42)
		src[4].SetInt(99)
		dst := NewHeapCells(2)
		CopyHeapCells(dst, src) // must not panic
		if got := dst[0].GetInt(); got != 42 {
			t.Errorf("dst[0] = %d, want 42", got)
		}
	})
}

// TestToSlot verifies descriptor-specific ToSlot conversion.
func TestToSlot(t *testing.T) {
	t.Run("int to Slot", func(t *testing.T) {
		c := &HeapCell{}
		c.SetInt(42)
		s := c.ToSlot("I")
		if s.Num() != 42 {
			t.Errorf("ToSlot(I).Num() = %d, want 42", s.Num())
		}
	})

	t.Run("float to Slot preserves bits", func(t *testing.T) {
		c := &HeapCell{}
		c.SetFloat(3.14)
		s := c.ToSlot("F")
		got := math.Float32frombits(uint32(s.Num()))
		if got != 3.14 {
			t.Errorf("ToSlot(F) round-trip = %v, want 3.14", got)
		}
	})

	t.Run("ref to Slot", func(t *testing.T) {
		c := &HeapCell{}
		obj := &Object{}
		c.SetRef(obj)
		s := c.ToSlot("Ljava/lang/Object;")
		if s.Ref() != obj {
			t.Error("ToSlot(L) ref not preserved")
		}
	})

	t.Run("long panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("ToSlot(J) should panic")
			}
		}()
		c := &HeapCell{}
		c.ToSlot("J")
	})

	t.Run("double panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("ToSlot(D) should panic")
			}
		}()
		c := &HeapCell{}
		c.ToSlot("D")
	})
}
