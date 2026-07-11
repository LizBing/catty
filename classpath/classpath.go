package classpath

// Classpath locates user classes. The Java core library (java.lang.Object,
// String, System, ...) is NOT loaded from a JRE here: catty implements those
// essential classes natively in the native package, so MVP needs no rt.jar.
type Classpath struct {
	user Entry
}

// Parse builds a Classpath from the -cp/-classpath option string. An empty
// option defaults to the current directory, matching java's behavior.
func Parse(cpOption string) *Classpath {
	if cpOption == "" {
		cpOption = "."
	}
	return &Classpath{user: newCompositeEntry(cpOption)}
}

func (c *Classpath) ReadClass(name string) ([]byte, Entry, error) {
	return c.user.ReadClass(name)
}

func (c *Classpath) String() string { return c.user.String() }
