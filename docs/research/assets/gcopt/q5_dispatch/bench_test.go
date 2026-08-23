package q5_dispatch

import "testing"

var sink int64

func benchTypeSwitch(b *testing.B, vals []I) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := vals[i&(len(vals)-1)]
		switch c := v.(type) {
		case *T1:
			sink += c.Val()
		case *T2:
			sink += c.Val()
		case *T3:
			sink += c.Val()
		case *T4:
			sink += c.Val()
		case *T5:
			sink += c.Val()
		case *T6:
			sink += c.Val()
		case *T7:
			sink += c.Val()
		case *T8:
			sink += c.Val()
		}
	}
}

// benchTypeSwitchNoBind deliberately calls the method on the interface value
// instead of the bound concrete, showing the double-dispatch pitfall the
// emitter must avoid.
func benchTypeSwitchNoBind(b *testing.B, vals []I) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := vals[i&(len(vals)-1)]
		switch v.(type) {
		case *T1:
			sink += v.Val()
		case *T2:
			sink += v.Val()
		case *T3:
			sink += v.Val()
		case *T4:
			sink += v.Val()
		case *T5:
			sink += v.Val()
		case *T6:
			sink += v.Val()
		case *T7:
			sink += v.Val()
		case *T8:
			sink += v.Val()
		}
	}
}

func benchAssertChain(b *testing.B, vals []I) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := vals[i&(len(vals)-1)]
		if c, ok := v.(*T1); ok {
			sink += c.Val()
		} else if c, ok := v.(*T2); ok {
			sink += c.Val()
		} else if c, ok := v.(*T3); ok {
			sink += c.Val()
		} else if c, ok := v.(*T4); ok {
			sink += c.Val()
		} else if c, ok := v.(*T5); ok {
			sink += c.Val()
		} else if c, ok := v.(*T6); ok {
			sink += c.Val()
		} else if c, ok := v.(*T7); ok {
			sink += c.Val()
		} else if c, ok := v.(*T8); ok {
			sink += c.Val()
		}
	}
}

func benchTagSwitch(b *testing.B, objs []*Obj) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := objs[i&(len(objs)-1)]
		switch o.tag {
		case 1:
			sink += o.n + 1
		case 2:
			sink += o.n + 2
		case 3:
			sink += o.n + 3
		case 4:
			sink += o.n + 4
		case 5:
			sink += o.n + 5
		case 6:
			sink += o.n + 6
		case 7:
			sink += o.n + 7
		case 8:
			sink += o.n + 8
		}
	}
}

// 2-way dispatch.
func BenchmarkTypeSwitch2(b *testing.B) {
	benchTypeSwitch(b, []I{newI1(), newI2()})
}

func BenchmarkAssertChain2(b *testing.B) {
	benchAssertChain(b, []I{newI1(), newI2()})
}

// 4-way dispatch.
func BenchmarkTypeSwitch4(b *testing.B) {
	benchTypeSwitch(b, []I{newI1(), newI2(), newI3(), newI4()})
}

func BenchmarkAssertChain4(b *testing.B) {
	benchAssertChain(b, []I{newI1(), newI2(), newI3(), newI4()})
}

// 8-way dispatch.
func BenchmarkTypeSwitch8(b *testing.B) {
	benchTypeSwitch(b, []I{newI1(), newI2(), newI3(), newI4(), newI5(), newI6(), newI7(), newI8()})
}

func BenchmarkTypeSwitchNoBind8(b *testing.B) {
	benchTypeSwitchNoBind(b, []I{newI1(), newI2(), newI3(), newI4(), newI5(), newI6(), newI7(), newI8()})
}

func BenchmarkAssertChain8(b *testing.B) {
	benchAssertChain(b, []I{newI1(), newI2(), newI3(), newI4(), newI5(), newI6(), newI7(), newI8()})
}

// Int-tag switch baseline (unified representation).
func BenchmarkTagSwitch8(b *testing.B) {
	objs := []*Obj{{1, 10}, {2, 20}, {3, 30}, {4, 40}, {5, 50}, {6, 60}, {7, 70}, {8, 80}}
	benchTagSwitch(b, objs)
}
