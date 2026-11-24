package rt

import "sync"

// Gate 5 extends the Gate 4 bootstrap inventory with the application lifecycle
// surface: java.lang.System.exit/nanoTime, the java.lang.Runtime singleton and
// its shutdown operations, and java.lang.IllegalStateException. The inventory
// is assembled once on first use on top of the Gate 4 definitions, so earlier
// profiles never publish these members.

var (
	gate5Once      sync.Once
	gate5Inventory []*bsClass
	gate5Index     map[string]*bsClass
)

func gate5BootstrapInventory() []*bsClass {
	gate5Once.Do(func() {
		gate5Inventory = buildGate5Inventory()
		gate5Index = map[string]*bsClass{}
		for _, c := range gate5Inventory {
			gate5Index[c.name] = c
		}
	})
	return gate5Inventory
}

func gate5BootstrapByName() map[string]*bsClass {
	gate5BootstrapInventory()
	return gate5Index
}

func buildGate5Inventory() []*bsClass {
	source := gate4BootstrapInventory()
	out := make([]*bsClass, 0, len(source)+2)
	merged := map[string]*bsClass{}
	for _, c := range source {
		nc := cloneBootstrapClass(c)
		out = append(out, nc)
		merged[nc.name] = nc
	}
	if system := merged["java/lang/System"]; system != nil {
		system.methods = append(system.methods,
			bsMethod{"exit", "(I)V", AccPublic | AccStatic, intrinsicSystemExit},
			bsMethod{"nanoTime", "()J", AccPublic | AccStatic, intrinsicSystemNanoTime},
		)
	}
	merged["java/lang/Runtime"] = &bsClass{
		name:  "java/lang/Runtime",
		super: "java/lang/Object",
		flags: AccPublic | AccFinal,
		methods: []bsMethod{
			{"getRuntime", "()Ljava/lang/Runtime;", AccPublic | AccStatic, intrinsicRuntimeGetRuntime},
			{"exit", "(I)V", AccPublic, intrinsicRuntimeExit},
			{"halt", "(I)V", AccPublic, intrinsicRuntimeHalt},
			{"addShutdownHook", "(Ljava/lang/Thread;)V", AccPublic, intrinsicRuntimeAddShutdownHook},
			{"removeShutdownHook", "(Ljava/lang/Thread;)Z", AccPublic, intrinsicRuntimeRemoveShutdownHook},
			{"gc", "()V", AccPublic, intrinsicRuntimeGC},
		},
	}
	merged["java/lang/IllegalStateException"] = &bsClass{
		name: "java/lang/IllegalStateException", super: "java/lang/RuntimeException",
		methods: throwableInit(),
	}
	// Keep declaration order deterministic: the Gate 4 inventory order followed
	// by the new classes in lexical order.
	out = append(out, merged["java/lang/IllegalStateException"], merged["java/lang/Runtime"])
	return out
}
