// Fixture for the stack-backfill acceptance test (DEBT-0019 diagnostic
// infrastructure). main -> stepA performs an astore on a null arrayref:
// the bootstrap NPE must carry a two-frame Java trace (stepA, main).
public class TraceProbe {
    public static void main(String[] args) {
        stepA(args.length);
    }

    static void stepA(int n) {
        int[] a = null;
        a[0] = n;
    }
}
