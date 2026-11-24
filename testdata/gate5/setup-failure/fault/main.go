// Command fault is a Go subprocess fixture that injects a host panic while a
// Java frame is active, to prove the Java frame snapshot and the original Go
// panic stack are both reported.
package main

import (
	"os"

	"catty/rt"
)

func main() {
	image := &rt.Image{
		ABI:       rt.Gate5ABI,
		Profile:   rt.Gate5ProfileID,
		Plugins:   true,
		Threads:   true,
		Lifecycle: true,
		Limits:    rt.Gate3PluginLimits,
		Classes: []*rt.ClassDef{{
			Name:  "Fault",
			Super: "java/lang/Object",
			Flags: rt.AccPublic,
			Methods: []*rt.MethodDef{{
				Owner:      "Fault",
				Name:       "main",
				Descriptor: "([Ljava/lang/String;)V",
				Flags:      rt.AccPublic | rt.AccStatic,
				Engine:     "aot",
				EntryIndex: -1,
				Entry: func(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
					system, err := t.RT.Class("java/lang/System")
					if err != nil {
						panic(err)
					}
					return t.Call(system.MethodIndex["nanoTime()J"], nil, nil)
				},
			}},
		}},
		Main: rt.MethodKey{Owner: "Fault", Name: "main", Descriptor: "([Ljava/lang/String;)V"},
	}
	status := rt.RunAppConfigured(image, os.Args[1:], func(r *rt.Runtime) {
		r.HostFaultHook = func(t *rt.Thread, m *rt.Method) {
			if m.Name == "nanoTime" {
				panic("injected host fault")
			}
		}
	})
	os.Exit(status)
}
