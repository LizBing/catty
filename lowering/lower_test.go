package lowering

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"catty/classloader"
	"catty/classpath"
	"catty/opcode"
)

// compileFixtures compiles every tests/fixtures/*.java into a temp dir and
// returns it, so a classloader can read the resulting .class files.
func compileFixtures(t *testing.T) string {
	t.Helper()
	src, err := filepath.Abs(filepath.Join("..", "tests", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-source", "8", "-target", "8", "-nowarn", "-d", out}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".java" {
			args = append(args, filepath.Join(src, e.Name()))
		}
	}
	if err := exec.Command("javac", args...).Run(); err != nil {
		t.Fatalf("javac failed: %v", err)
	}
	return out
}

// TestLowerAllFixtures is the A0.1 soundness gate: every interpreted method of
// every fixture must lower with a consistent depth dataflow (no pc reached at
// two depths), and every instruction's Uses/Defs must stay in range.
func TestLowerAllFixtures(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	classes := []string{
		"HelloWorld", "Fibonacci", "Factorial", "ArraySum",
		"OOPDemo", "Point", "StaticFields", "SwitchDemo",
	}

	lowered := 0
	for _, name := range classes {
		c := cl.LoadClass(name)
		for _, m := range c.Methods() {
			if m.IsNative() || len(m.Code()) == 0 {
				continue
			}
			ir, err := Lower(m)
			if err != nil {
				t.Errorf("lower %s.%s%s: %v", name, m.Name(), m.Descriptor(), err)
				continue
			}
			maxStack := int(m.MaxStack())
			for pc := 0; pc < len(ir.Insts); pc++ {
				inst := &ir.Insts[pc]
				if !inst.Present {
					continue
				}
				for _, u := range inst.Uses {
					if int(u) < 0 || int(u) >= maxStack {
						t.Errorf("%s.%s pc=%d use slot %d out of range [0,%d)",
							name, m.Name(), pc, u, maxStack)
					}
				}
				for _, d := range inst.Defs {
					if int(d) < 0 || int(d) >= maxStack {
						t.Errorf("%s.%s pc=%d def slot %d out of range [0,%d)",
							name, m.Name(), pc, d, maxStack)
					}
				}
			}
			lowered++
		}
	}
	if lowered == 0 {
		t.Fatal("no methods were lowered")
	}
	t.Logf("lowered %d methods across %d classes without dataflow errors", lowered, len(classes))
}

// TestLowerHelloWorldInit hand-checks the eliminated stack on a tiny known
// method. HelloWorld.<init>()V is roughly:
//
//	0: aload_0            ; depth 0 → defs [0]
//	1: invokespecial …    ; pop 1 (this), depth 1 → uses [0]
//	4: return             ; depth 0
//
// so the IR should show the constructor's `this` flowing from aload_0 into the
// invokespecial via slot 0, with no operand stack surviving to the return.
func TestLowerHelloWorldInit(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	c := cl.LoadClass("HelloWorld")
	init := c.GetMethod("<init>", "()V")
	if init == nil {
		t.Fatal("<init> not found")
	}
	ir, err := Lower(init)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}

	// pc 0 is aload_0: defines slot 0 at entry depth 0.
	a0 := ir.Insts[0]
	if a0.Op != opcode.Aload0 {
		t.Fatalf("pc0 op = %v, want aload_0", a0.Op)
	}
	if len(a0.Defs) != 1 || a0.Defs[0] != 0 {
		t.Errorf("aload_0 defs = %v, want [0]", a0.Defs)
	}

	// The invokespecial consumes `this` from slot 0.
	var spec *IRInst
	for i := range ir.Insts {
		if ir.Insts[i].Present && ir.Insts[i].Op == opcode.Invokespecial {
			spec = &ir.Insts[i]
			break
		}
	}
	if spec == nil {
		t.Fatal("no invokespecial found in <init>")
	}
	if len(spec.Uses) != 1 || spec.Uses[0] != 0 {
		t.Errorf("invokespecial uses = %v, want [0] (this)", spec.Uses)
	}
	if len(spec.Defs) != 0 {
		t.Errorf("invokespecial defs = %v, want []", spec.Defs)
	}
}
