// Compiles alone to a Dep.class that has neither good() nor value. The
// acceptance harness runs Link against a classpath of the main classes (with
// Dep.class removed) followed by this variant directory.
final class Dep {
    static int other() {
        return 1;
    }
}
