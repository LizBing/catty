// Category: class initialization — predecessor failure.
// A failing superclass makes the subclass erroneous without running SubFail's
// <clinit>. Expected: EIIE, NCDFE, and no "SubFail" line.
class SuperFail {
    static int x = fail();
    static int fail() {
        throw new RuntimeException("boom");
    }
}
class SubFail extends SuperFail {
    static {
        System.out.println("SubFail");
    }
}
public class SuperclassInitializationFailure {
    public static void main(String[] args) {
        try {
            new SubFail();
        } catch (Throwable value) {
            System.out.println(value.getClass().getName());
        }
        try {
            new SubFail();
        } catch (Throwable value) {
            System.out.println(value.getClass().getName());
        }
    }
}
