public class Bench {
    static int now() { return System.tickMillis(); }

    // 1. integer arithmetic loop
    static int arith(int n) {
        int t0 = now();
        int s = 0;
        for (int i = 0; i < n; i++) s += i * 3 - (i >> 1);
        int t1 = now();
        if (s == 42) System.out.println("impossible");
        return t1 - t0;
    }

    // 2. virtual dispatch: monomorphic-ish call in hot loop
    static class Node { int v; int get() { return v; } Node(int v) { this.v = v; } }
    static int vcall(int n) {
        int t0 = now();
        Node nd = new Node(7);
        int s = 0;
        for (int i = 0; i < n; i++) s += nd.get();
        int t1 = now();
        if (s == 42) System.out.println("impossible");
        return t1 - t0;
    }

    // 3. HashMap put/get mix
    static int mapops(int n) {
        java.util.HashMap m = new java.util.HashMap();
        int t0 = now();
        for (int i = 0; i < n; i++) {
            String k = "k" + (i % 1024);
            Integer c = (Integer) m.get(k);
            m.put(k, Integer.valueOf(c == null ? 1 : c.intValue() + 1));
        }
        int t1 = now();
        return t1 - t0;
    }

    // 4. string concat via StringBuilder chain
    static int strcat(int n) {
        int t0 = now();
        String acc = "";
        for (int i = 0; i < n; i++) acc = new StringBuilder(acc).append("x").append(i % 10).toString();
        int t1 = now();
        return t1 - t0;
    }

    public static void main(String[] args) {
        int N = Integer.parseInt(args[0]);
        int W = N / 5;
        // JIT warmup (no-op for Catty)
        arith(W); vcall(W); mapops(W); 
        for (int r = 0; r < 3; r++) {
            System.out.println("arith,"   + arith(N));
            System.out.println("vcall,"   + vcall(N));
            System.out.println("mapops,"  + mapops(N));
            System.out.println("strcat,"  + strcat(N / 10));
            System.out.println("---");
        }
    }
}
