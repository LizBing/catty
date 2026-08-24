// R-0008 spike B: interop boundary tax. Two tight loops, one calling a
// Go-registered native, one calling an equivalent pure-Java static —
// same shape, same loop overhead; the delta is the boundary tax.
public class GoCallBench {
    static int javaAdd(int a, int b) { return a + b; }

    static long timed(int n, boolean useGo) {
        long t0 = System.nanoTime();
        int s = 0;
        for (int i = 0; i < n; i++) {
            if (useGo) s += go.Bridge.add(i, 1);
            else s += javaAdd(i, 1);
        }
        long t1 = System.nanoTime();
        if (s == -1) System.out.println("impossible");
        return t1 - t0;
    }

    public static void main(String[] args) {
        int n = Integer.parseInt(args[0]);
        // warmup both paths
        for (int r = 0; r < 3; r++) { timed(n / 5, true); timed(n / 5, false); }
        for (int r = 0; r < 3; r++) {
            System.out.println("go_ns," + timed(n, true));
            System.out.println("java_ns," + timed(n, false));
            System.out.println("---");
        }
    }
}
