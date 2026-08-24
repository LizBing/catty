// P-0009 U4: sustained-load latency percentiles. Each iteration is one
// "request": map get/put + string concat + virtual call. Latencies land in
// a preallocated long[]; percentiles come from an in-Java heapsort so no
// Arrays.sort surface is needed.
public class P99Bench {
    static class Node { int v; int get() { return v; } Node(int v) { this.v = v; } }

    // heapsort on arr[0..n)
    static void sort(long[] a, int n) {
        for (int i = n / 2 - 1; i >= 0; i--) sift(a, i, n);
        for (int i = n - 1; i > 0; i--) {
            long t = a[0]; a[0] = a[i]; a[i] = t;
            sift(a, 0, i);
        }
    }

    static void sift(long[] a, int root, int n) {
        while (2 * root + 1 < n) {
            int c = 2 * root + 1;
            if (c + 1 < n && a[c + 1] > a[c]) c++;
            if (a[root] >= a[c]) return;
            long t = a[root]; a[root] = a[c]; a[c] = t;
            root = c;
        }
    }

    static long pct(long[] sorted, int n, double p) {
        int idx = (int) Math.floor(p * (n - 1));
        return sorted[idx];
    }

    public static void main(String[] args) throws Exception {
        int warm = Integer.parseInt(args[0]);
        int measured = Integer.parseInt(args[1]);
        java.util.HashMap m = new java.util.HashMap();
        Node nd = new Node(7);

        long[] lat = new long[measured];
        for (int i = 0; i < warm; i++) {
            String k = "k" + (i % 1024);
            Integer c = (Integer) m.get(k);
            m.put(k, Integer.valueOf(c == null ? 1 : c.intValue() + 1));
            StringBuilder sb = new StringBuilder(k).append(':').append(nd.get());
            if (sb.toString().length() == 99) System.out.println("impossible");
        }
        for (int i = 0; i < measured; i++) {
            long t0 = System.nanoTime();
            String k = "k" + (i % 1024);
            Integer c = (Integer) m.get(k);
            m.put(k, Integer.valueOf(c == null ? 1 : c.intValue() + 1));
            StringBuilder sb = new StringBuilder(k).append(':').append(nd.get());
            if (sb.toString().length() == 99) System.out.println("impossible");
            lat[i] = System.nanoTime() - t0;
        }
        sort(lat, measured);
        System.out.println("p50_ns=" + pct(lat, measured, 0.50));
        System.out.println("p90_ns=" + pct(lat, measured, 0.90));
        System.out.println("p99_ns=" + pct(lat, measured, 0.99));
        System.out.println("max_ns=" + lat[measured - 1]);
    }
}
