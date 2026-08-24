// Single-load fixture: monomorphic virtual call in hot loop (P-0009 T0.5).
public class VCallBench {
    static class Node { int v; int get() { return v; } Node(int v) { this.v = v; } }

    public static void main(String[] args) {
        int n = Integer.parseInt(args[0]);
        Node nd = new Node(7);
        long t0 = System.nanoTime();
        int s = 0;
        for (int i = 0; i < n; i++) s += nd.get();
        long t1 = System.nanoTime();
        if (s == 42) System.out.println("impossible");
        System.out.println("vcall_ns," + (t1 - t0));
    }
}
