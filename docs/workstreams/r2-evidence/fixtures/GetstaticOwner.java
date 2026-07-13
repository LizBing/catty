// Category: class initialization trigger — getstatic, declarer rule (JVMS §5.5).
// Probes: a static field is referenced through a subclass (Sub.x), but the
// field is declared in Base. Only the declaring class (Base) is initialized;
// the subclass (Sub) is NOT initialized merely because the field was reached
// through it.
// Expected (Temurin 25): "Base" then "7"; Sub.<clinit> does NOT run.
class Base {
    static int x;
    static {
        System.out.println("Base");
        x = 7;
    }
}
class Sub extends Base {
    static {
        System.out.println("Sub");
    }
}
public class GetstaticOwner {
    public static void main(String[] args) {
        System.out.println(Sub.x);
    }
}
