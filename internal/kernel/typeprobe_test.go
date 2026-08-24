package kernel

import "testing"

func TestIntegerTypeProbe(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	c, ok := k.ClassByName("java/lang/Integer")
	if !ok { t.Fatal("no Integer") }
	t.Logf("fields=%d declared=%d statics=%d", len(c.fieldsByKey), len(c.DeclaredFields), len(c.Statics))
	f := c.fieldsByKey[memberKey("TYPE", "Ljava/lang/Class;")]
	t.Logf("TYPE field=%v", f)
	if f != nil {
		t.Logf("static slot=%d value=%v", f.StaticSlot, c.Statics[f.StaticSlot])
	}
}
