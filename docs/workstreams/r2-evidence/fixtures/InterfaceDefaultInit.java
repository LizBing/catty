// Category: INTERFACE initialization during implementing-class initialization.
// Probes: initializing a class must first initialize a superinterface that
// declares a default method. invokeinterface itself is NOT an init trigger.
// Expected (Temurin 25): "I", "after-new", then "I.m".
// Expected (catty R1): only "I.m" — interpreter/invoke.go skips interface init.
interface Iface {
    int X = mark();                       // forces Iface.<clinit> to exist
    static int mark() {
        System.out.println("I");
        return 1;
    }
    default void m() {
        System.out.println("I.m");
    }
}
class Impl implements Iface {}
public class InterfaceDefaultInit {
    public static void main(String[] args) {
        Impl value = new Impl();
        System.out.println("after-new");
        value.m();
    }
}
