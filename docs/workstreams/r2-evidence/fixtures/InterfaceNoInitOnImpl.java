// Category: INTERFACE initialization — no init on mere implementation (JVMS §5.5).
// Probes: creating an instance of a class that implements an interface does NOT
// initialize the interface (no default method is invoked, no static field read).
// Expected (Temurin 25 and catty): "main" only; Iface2.<clinit> does NOT run.
interface Iface2 {
    int X = mark2();
    static int mark2() {
        System.out.println("I2");
        return 1;
    }
}
class Impl2 implements Iface2 {}
public class InterfaceNoInitOnImpl {
    public static void main(String[] args) {
        new Impl2();
        System.out.println("main");
    }
}
