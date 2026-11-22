// Deferred-linkage fixture: Link.java compiles against the real Dep, but the
// acceptance run deletes Dep.class from the main classes and resolves against
// variant/DepChanged.java, whose Dep has neither good() nor value.
//
// Startup must succeed with the unused bad references present. Each use site
// then resolves lazily and fails with its specified linkage error; the second
// call shows the failure is repeatable and leaves the class usable.
public final class Link {
    static int startup() {
        return 1;
    }

    static int finish() {
        return 5;
    }

    static void unusedCall() {
        Dep.good();
    }

    static void unusedField() {
        int ignored = Dep.value;
    }

    static int callGood() {
        return Dep.good() + 1;
    }

    static int readValue() {
        return Dep.value + 1;
    }

    public static void main(String[] args) {
        System.out.println(startup());
        try {
            System.out.println(callGood());
        } catch (NoSuchMethodError expected) {
            System.out.println(2);
        }
        try {
            System.out.println(callGood());
        } catch (NoSuchMethodError expected) {
            System.out.println(3);
        }
        try {
            System.out.println(readValue());
        } catch (NoSuchFieldError expected) {
            System.out.println(4);
        }
        System.out.println(finish());
    }
}
