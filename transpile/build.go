package transpile

import (
	"fmt"
	"strings"

	"catty/classloader"
	"catty/classpath"
	"catty/lowering"
	"catty/opcode"
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

	// Conservative build-time rejection: AOT exception propagation is not
	// yet wired for class-initialization failure. Any emitted method whose
	// IR contains an init-triggering instruction targeting a class with a
	// <clinit> is rejected — whether or not it has exception handlers —
	// because neither caught nor uncaught EIIE/NCDFE can be reported
	// without a Go panic stack trace. Cross-engine exception propagation
	// belongs to a separate future workstream.
	for _, p := range emittable {
		if reason := hasInitTrigger(p.method, p.ir, cl); reason != "" {
			return "", fmt.Errorf("transpile: %s.%s%s triggers class initialization on %s; AOT exception propagation not yet supported",
				p.method.Owner().Name(), p.method.Name(), p.method.Descriptor(), reason)
		}
	}

	// Conservative build-time rejection: AOT concurrency execution is not
	// implemented (ADR-0028, ADR-0029). Reject the build if any method uses
	// monitor bytecodes or synchronized flag, or if the program touches
	// java/lang/Thread (the execution-context ABI is not wired).
	for _, cls := range cl.Classes() {
		for _, m := range cls.Methods() {
			if reason := concurrencyReason(m); reason != "" {
				return "", fmt.Errorf("transpile: %s.%s%s — %s; AOT concurrency not supported",
					cls.Name(), m.Name(), m.Descriptor(), reason)
			}
		}
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

	// Assemble: emitted funcs + simple main() wrapper.
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

// hasInitTrigger scans a method's IR for instructions that trigger class
// initialization (getstatic, putstatic, new, invokestatic) and computes the
// full initialization predecessor closure (ADR-0025 / JVMS §5.5) for each
// target. If any class or interface in the closure has a <clinit>, the
// entire build is rejected as Not implemented.
//
// The closure includes the target class itself, its superclass chain, and
// for each class in that chain, its recursively-enumerated default-bearing
// superinterfaces. This is the same set that InitializeClass would traverse
// at runtime; any member of the closure could fail during initialization,
// and without cross-engine exception propagation a Go panic would result.
func hasInitTrigger(method *rtda.Method, ir *lowering.IR, loader rtda.Loader) string {
	cp := method.Owner().ConstantPool()
	if cp == nil {
		return ""
	}
	for _, inst := range ir.Insts {
		if !inst.Present {
			continue
		}
		var target *rtda.Class
		switch inst.Op {
		case opcode.Getstatic, opcode.Putstatic:
			refClass, name, desc := cp.MemberRef(inst.Index)
			if c := loader.LoadClass(refClass); c != nil {
				if field := c.LookupField(name, desc); field != nil {
					target = field.Owner()
				}
			}
		case opcode.Invokestatic:
			refClass, name, desc := cp.MemberRef(inst.Index)
			if c := loader.LoadClass(refClass); c != nil {
				if m := c.LookupMethod(name, desc); m != nil {
					target = m.Owner()
				}
			}
		case opcode.New:
			className := cp.ClassName(inst.Index)
			target = loader.LoadClass(className)
		}
		if target != nil {
			if reason := initClosureHasClinit(target); reason != "" {
				return reason
			}
		}
	}
	return ""
}

// initClosureHasClinit checks whether any class or interface in the
// initialization predecessor closure of `class` (per ADR-0025 / JVMS §5.5)
// defines a <clinit>. Returns the name of the first match, or "" if none.
func initClosureHasClinit(class *rtda.Class) string {
	seen := make(map[string]bool)
	return checkClosure(class, seen)
}

func checkClosure(class *rtda.Class, seen map[string]bool) string {
	if class == nil || seen[class.Name()] {
		return ""
	}
	seen[class.Name()] = true

	if class.GetMethod("<clinit>", "()V") != nil {
		return class.Name()
	}

	// For classes (not interfaces): recurse into superclass and
	// default-bearing superinterfaces per JVMS §5.5 step 7.
	if !class.IsInterface() {
		// Superclass chain.
		if class.SuperClass() != nil {
			if reason := checkClosure(class.SuperClass(), seen); reason != "" {
				return reason
			}
		}
		// Default-bearing superinterfaces (depth-first, left-to-right).
		for _, iface := range class.DefaultBearingSuperInterfaces(make(map[string]bool)) {
			if reason := checkClosure(iface, seen); reason != "" {
				return reason
			}
		}
	}

	return ""
}
