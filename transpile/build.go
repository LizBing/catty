package transpile

import (
	"fmt"
	"strings"

	"catty/classloader"
	"catty/classpath"
	"catty/lowering"
	"catty/rtda"
)

// BuildProgram transpiles a whole program: loads the main class and all
// reachable classes, AOT-emits every emittable method, and wraps them in a
// standalone Go main package that Bootstraps the runtime and calls the
// transpiled main. The result is compiled with `go build` to produce a native
// binary.
//
// Methods the emitter can't handle (unsupported opcodes, instance methods) are
// skipped; at runtime they're served by the interpreter via the bridge.
func BuildProgram(mainClass, classpathStr string) (string, error) {
	cl := classloader.New(classpath.Parse(classpathStr))
	loadReachable(cl, mainClass)

	// Pass 1: try to emit each real bytecode method; build the emittable set.
	type pending struct {
		method *rtda.Method
		ir     *lowering.IR
	}
	var emittable []pending
	emitted := map[string]bool{}
	for _, cls := range cl.Classes() {
		for _, m := range cls.Methods() {
			if m.IsNative() || len(m.Code()) == 0 {
				continue
			}
			ir, err := lowering.Lower(m)
			if err != nil {
				continue
			}
			if _, err := Emit(m, ir, cl, nil); err != nil {
				continue
			}
			key := cls.Name() + "\x00" + m.Name() + "\x00" + m.Descriptor()
			emitted[key] = true
			emittable = append(emittable, pending{m, ir})
		}
	}

	// Check main is emittable.
	mainCls := cl.LoadClass(mainClass)
	mainMethod := mainCls.GetMethod("main", "([Ljava/lang/String;)V")
	if mainMethod == nil {
		return "", fmt.Errorf("transpile: main method not found in %s", mainClass)
	}
	mainKey := mainClass + "\x00main\x00([Ljava/lang/String;)V"
	if !emitted[mainKey] {
		return "", fmt.Errorf("transpile: %s.main cannot be AOT'd (unsupported opcodes?)", mainClass)
	}

	// Pass 2: re-emit with the emittable set (invokestatic dispatch).
	var src strings.Builder
	for _, p := range emittable {
		s, err := Emit(p.method, p.ir, cl, emitted)
		if err != nil {
			return "", fmt.Errorf("transpile: re-emit %s.%s: %v", p.method.Owner().Name(), p.method.Name(), err)
		}
		src.WriteString(s)
		src.WriteString("\n")
	}

	// Assemble: emitted funcs + main() wrapper.
	program := "package main\n\nimport (\n\t\"catty/runtime\"\n\t\"catty/rtda\"\n)\n\n" +
		src.String() +
		"\nfunc main() {\n\truntime.Bootstrap(" + fmt.Sprintf("%q, %q", classpathStr, mainClass) + ")\n\t" +
		mangle(mainClass, "main") + "((*rtda.Object)(nil))\n}\n"
	return program, nil
}

// loadReachable loads the main class and all classes transitively referenced
// via constant pool CONSTANT_Class entries (the reachability closure).
func loadReachable(cl *classloader.ClassLoader, mainClass string) {
	visited := map[string]bool{}
	worklist := []string{mainClass}
	for len(worklist) > 0 {
		name := worklist[0]
		worklist = worklist[1:]
		if visited[name] || strings.HasPrefix(name, "[") {
			continue
		}
		visited[name] = true
		cls := cl.LoadClass(name)
		cp := cls.ConstantPool()
		if cp == nil {
			continue // synthetic (native) class — no constant pool
		}
		for _, refName := range cp.ClassRefNames() {
			if !visited[refName] && !strings.HasPrefix(refName, "[") {
				worklist = append(worklist, refName)
			}
		}
	}
}
