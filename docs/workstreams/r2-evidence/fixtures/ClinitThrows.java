// Category: class initialization FAILURE semantics (JVMS §5.5; JLS §12.4.2).
// Probes: <clinit> throws. First triggering access must surface
// ExceptionInInitializerError; the class becomes "erroneous" and a SECOND
// access must surface NoClassDefFoundError.
// Expected (Temurin 25):
//   java.lang.ExceptionInInitializerError
//   java.lang.NoClassDefFoundError
// Expected (catty R1): the raw RuntimeException propagates (no EIIE wrapping,
// no erroneous state); behavior on second access depends on R1 initStarted.
class Bomb {
    static int x = 5;
    static { boom(); }                       // indirection so the initializer
    static void boom() {                     // can complete normally (JLS).
        throw new RuntimeException("boom");
    }
}
public class ClinitThrows {
    public static void main(String[] args) {
        try {
            System.out.println(Bomb.x);
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
        try {
            System.out.println(Bomb.x);
        } catch (Throwable t) {
            System.out.println(t.getClass().getName());
        }
    }
}
