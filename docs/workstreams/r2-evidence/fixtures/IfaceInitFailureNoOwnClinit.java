// Category: class initialization — predecessor failure through
// default-bearing superinterface chain when the target class has no
// <clinit> of its own.
// Probes: new ImplNoClinit triggers DefIface.<clinit> (a default-bearing
// interface) which fails; ImplNoClinit has no <clinit> itself — the
// failure comes solely from the required interface predecessor. This
// proves the init predecessor closure must include default-bearing
// superinterfaces, not just the superclass chain.
// Expected (Temurin 25 and catty): EIIE then NCDFE.
interface DefIface {
    int X = fail();
    static int fail() {
        throw new RuntimeException("iface init failed");
    }
    default void m() {}
}
class ImplNoClinit implements DefIface {
    // No <clinit> — init failure comes solely from the default-bearing interface.
}
public class IfaceInitFailureNoOwnClinit {
    public static void main(String[] args) {
        try {
            new ImplNoClinit();
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
        try {
            new ImplNoClinit();
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
    }
}
