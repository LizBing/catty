// Category: class initialization — predecessor failure through superclass chain
// when the target class has no <clinit> of its own.
// Probes: new ChildNoClinit triggers ParentWithClinit.<clinit> which fails;
// ChildNoClinit has no <clinit> itself — the failure comes solely from the
// superclass chain. This proves the init predecessor closure must be checked
// at AOT build time, not just the target class's own <clinit>.
// Expected (Temurin 25 and catty): EIIE then NCDFE.
class ParentWithClinit {
    static int x = fail();
    static int fail() {
        throw new RuntimeException("parent init failed");
    }
}
class ChildNoClinit extends ParentWithClinit {
    // No <clinit> — init failure comes solely from the superclass.
}
public class SuperInitFailureNoOwnClinit {
    public static void main(String[] args) {
        try {
            new ChildNoClinit();
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
        try {
            new ChildNoClinit();
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
    }
}
