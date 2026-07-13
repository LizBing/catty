// Category: bootstrap resolution boundary (PROJECT_STATUS explicit boundary).
// Probes: Integer.toString reaches jdk.internal.misc.Unsafe on JDK 25.
// Expected (Temurin 25): "42".
// Expected (catty R1, pure-synthetic): class java/lang/Integer not served
// (no java.base); records the boundary as Not implemented, not silent approx.
public class ReachUnsafe {
    public static void main(String[] args) {
        System.out.println(Integer.toString(42));
    }
}
