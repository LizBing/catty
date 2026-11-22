package rt

// Gate 3 extends the closed Gate 2 bootstrap inventory with the plugin API,
// loader mirrors and the additional checked/unchecked exception types. All
// declarations are still fixed profile metadata; no application class may
// define anything in these namespaces.

var (
	gate3Inventory []*bsClass
	gate3ByName    map[string]*bsClass
)

func init() {
	gate3Inventory = buildGate3Inventory()
	gate3ByName = map[string]*bsClass{}
	for _, c := range gate3Inventory {
		gate3ByName[c.name] = c
	}
}

func cloneBootstrapClass(c *bsClass) *bsClass {
	out := &bsClass{
		name:       c.name,
		super:      c.super,
		interfaces: append([]string(nil), c.interfaces...),
		flags:      c.flags,
	}
	out.fields = append([]bsField(nil), c.fields...)
	out.methods = append([]bsMethod(nil), c.methods...)
	return out
}

func buildGate3Inventory() []*bsClass {
	out := make([]*bsClass, 0, len(bootstrapInventory)+8)
	merged := map[string]*bsClass{}
	for _, c := range bootstrapInventory {
		nc := cloneBootstrapClass(c)
		out = append(out, nc)
		merged[nc.name] = nc
	}
	if class := merged["java/lang/Class"]; class != nil {
		class.methods = append(class.methods,
			bsMethod{"forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", AccPublic | AccStatic, intrinsicClassForName},
			bsMethod{"getClassLoader", "()Ljava/lang/ClassLoader;", AccPublic, intrinsicClassGetClassLoader},
		)
	}
	merged["java/lang/ClassLoader"] = &bsClass{
		name:  "java/lang/ClassLoader",
		super: "java/lang/Object",
		flags: AccPublic | AccAbstract,
	}
	merged["java/lang/ClassNotFoundException"] = &bsClass{
		name: "java/lang/ClassNotFoundException", super: "java/lang/ReflectiveOperationException",
		methods: throwableInit(),
	}
	merged["java/io/IOException"] = &bsClass{
		name: "java/io/IOException", super: "java/lang/Exception",
		methods: throwableInit(),
	}
	merged["java/lang/SecurityException"] = &bsClass{
		name: "java/lang/SecurityException", super: "java/lang/RuntimeException",
		methods: throwableInit(),
	}
	merged["catty/runtime/Plugins"] = &bsClass{
		name:  "catty/runtime/Plugins",
		super: "java/lang/Object",
		flags: AccPublic | AccFinal,
		methods: []bsMethod{
			{"readJar", "(Ljava/lang/String;)[B", AccPublic | AccStatic, intrinsicPluginsReadJar},
			{"register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", AccPublic | AccStatic, intrinsicPluginsRegister},
		},
	}
	merged["catty/runtime/PluginAdmissionException"] = &bsClass{
		name:  "catty/runtime/PluginAdmissionException",
		super: "java/lang/RuntimeException",
		flags: AccPublic | AccFinal,
		methods: []bsMethod{
			{"getCode", "()Ljava/lang/String;", AccPublic, intrinsicPluginAdmissionGetCode},
		},
	}
	// Keep declaration order deterministic: Gate 2 inventory order followed by
	// the Gate 3 additions in lexical order.
	order := []string{
		"java/lang/ClassLoader",
		"java/lang/ClassNotFoundException",
		"java/io/IOException",
		"java/lang/SecurityException",
		"catty/runtime/Plugins",
		"catty/runtime/PluginAdmissionException",
	}
	seen := map[string]bool{}
	for _, c := range out {
		seen[c.name] = true
	}
	for _, name := range order {
		if seen[name] {
			continue
		}
		out = append(out, merged[name])
		seen[name] = true
	}
	return out
}

func gate3BootstrapInventory() []*bsClass { return gate3Inventory }

func gate3BootstrapByName() map[string]*bsClass { return gate3ByName }
