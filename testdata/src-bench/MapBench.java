// Single-load fixture: HashMap get/put + boxing + string-key concat only
// (P-0009 T0.5).
public class MapBench {
    public static void main(String[] args) {
        int n = Integer.parseInt(args[0]);
        java.util.HashMap m = new java.util.HashMap();
        long t0 = System.nanoTime();
        for (int i = 0; i < n; i++) {
            String k = "k" + (i % 1024);
            Integer c = (Integer) m.get(k);
            m.put(k, Integer.valueOf(c == null ? 1 : c.intValue() + 1));
        }
        long t1 = System.nanoTime();
        System.out.println("mapops_ns," + (t1 - t0));
    }
}
