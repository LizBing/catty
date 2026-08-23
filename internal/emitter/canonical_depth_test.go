package emitter

import "testing"

// TestNetStackEffectSentinels pins effect values for opcodes whose loss
// previously caused silent +1 slot drift across whole methods (pop and
// aaload each went missing during table refactors).
func TestNetStackEffectSentinels(t *testing.T) {
	cases := []struct {
		op   byte
		name string
		want int
	}{
		{0x57, "pop", -1},
		{0x32, "aaload", -1},
		{0x2e, "iaload", -1},
		{0x3a, "astore", -1},
		{0xbb, "new", 1},
		{0x59, "dup", 1},
		{0x37, "lstore", -2},
		{0x14, "ldc2_w", 2},
		{0xa7, "goto", 0},
	}
	for _, c := range cases {
		if got := netStackEffect(c.op); got != c.want {
			t.Errorf("%s (%#x): got %d, want %d", c.name, c.op, got, c.want)
		}
	}
}
