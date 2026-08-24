package classfile

import "testing"

// TestOpcodeStackEffectTotal pins the contract that every byte value is
// classified — no silent defaults. A new opcode reaching an engine without
// a row must fail HERE, not as runtime stack corruption (DEBT-0019 class).
func TestOpcodeStackEffectTotal(t *testing.T) {
	for op := 0; op <= 255; op++ {
		_, class := StackEffect(byte(op))
		switch class {
		case EffectFixed, EffectDescriptor, EffectIllegal:
			// classified
		default:
			t.Fatalf("opcode %#02x: unclassified effect %d", op, class)
		}
	}
}

// TestOpcodeStackEffectPins pins the exact rows whose loss historically
// produced engine drift (DEBT-0015 / DEBT-0017 / DEBT-0019).
func TestOpcodeStackEffectPins(t *testing.T) {
	pins := []struct {
		op   byte
		name string
		want int
	}{
		{0x53, "aastore", -3},          // DEBT-0019 root cause
		{0x4f, "iastore", -3},
		{0x50, "lastore", -4},
		{0x52, "dastore", -4},
		{0x54, "bastore", -3},
		{0x55, "castore", -3},
		{0x56, "sastore", -3},
		{0x57, "pop", -1},              // DEBT-0015 loss
		{0x2f, "laload", 0},
		{0x31, "daload", 0},
		{0x5c, "dup2", 2},
		{0x5d, "dup2_x1", 2},
		{0x5e, "dup2_x2", 2},
		{0x94, "lcmp", -3},
		{0x97, "dcmpl", -3},
		{0x98, "dcmpg", -3},
		{0xaa, "tableswitch", -1},
		{0xab, "lookupswitch", -1},
		{0x85, "i2l", 1},
		{0x88, "l2i", -1},
		{0x8a, "l2d", 0},
		{0x75, "lneg", 0},
		{0xad, "lreturn", -2},
	}
	for _, p := range pins {
		got, class := StackEffect(p.op)
		if class != EffectFixed || got != p.want {
			t.Errorf("%#02x %s = (%d, class %d), want (%d, EffectFixed)",
				p.op, p.name, got, class, p.want)
		}
	}

	// Descriptor-dependent families must not pretend to be fixed.
	for _, op := range []byte{0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xc4, 0xc5} {
		if _, class := StackEffect(op); class != EffectDescriptor {
			t.Errorf("op %#02x: want EffectDescriptor, got %d", op, class)
		}
	}

	// Outlawed shapes stay loud.
	for _, op := range []byte{0xa8, 0xa9, 0xba} {
		if _, class := StackEffect(op); class != EffectIllegal {
			t.Errorf("op %#02x: want EffectIllegal, got %d", op, class)
		}
	}
}
