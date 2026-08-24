// P-0009 T4 experiment: N threads doing equal chunks of mixed work
// (virtual calls + arithmetic + boxed map updates), wall-clock via
// nanoTime. First concurrency-throughput datapoint for selling point #3.
public class ConcurrentBench {
    static class Worker extends Thread {
        int chunk;
        long sink;

        Worker(int chunk) { this.chunk = chunk; }

        public void run() {
            Node nd = new Node(7);
            java.util.HashMap m = new java.util.HashMap();
            long acc = 0;
            for (int i = 0; i < chunk; i++) {
                acc += nd.get();
                if ((i & 1023) == 0) {
                    m.put(nd, Integer.valueOf(i));
                }
            }
            sink = acc;
        }
    }

    static class Node { int v; int get() { return v; } Node(int v) { this.v = v; } }

    public static void main(String[] args) throws Exception {
        int threads = Integer.parseInt(args[0]);
        int chunk = Integer.parseInt(args[1]);
        Worker[] ws = new Worker[threads];
        long t0 = System.nanoTime();
        for (int i = 0; i < threads; i++) {
            ws[i] = new Worker(chunk);
            ws[i].start();
        }
        for (int i = 0; i < threads; i++) {
            ws[i].join();
        }
        long t1 = System.nanoTime();
        long sum = 0;
        for (int i = 0; i < threads; i++) {
            sum += ws[i].sink;
        }
        System.out.println("workers=" + threads + " chunk=" + chunk +
            " wall_ns=" + (t1 - t0) + " sink=" + (sum == 42 ? "impossible" : "ok"));
    }
}
