package gate5bench;

/**
 * Memory and thread kernels:
 *
 * <ul>
 *   <li>escaping allocation into a bounded ring,</li>
 *   <li>a field- and array-heavy kernel run at one and at four platform
 *       threads over disjoint slots of the same shared array,</li>
 *   <li>a contended monitor kernel with a fixed number of platform threads and
 *       verified completed work.</li>
 * </ul>
 *
 * Every kernel's final values and checksum are deterministic: threads touch
 * disjoint array slots and per-thread cells; the monitor counter's final value
 * is exact because every increment is synchronized. Threads are started and
 * joined inside each timed batch, so the per-batch time for the threaded
 * kernels includes fixed thread start/join cost.
 */
final class Memory {
    private Memory() {
    }

    static int allocKernel(Escape[] ring, int ops, int seed) {
        int acc = seed;
        int i = 0;
        int mask = ring.length - 1;
        while (i < ops) {
            Escape escape = new Escape();
            escape.id = i + seed;
            escape.extra = (i * 3) ^ seed;
            escape.payload = new int[4];
            escape.payload[0] = i;
            escape.payload[1] = escape.id;
            escape.payload[2] = escape.id * 3 + 1;
            escape.payload[3] = escape.payload[0] ^ escape.payload[1];
            ring[i & mask] = escape;
            acc = acc + escape.payload[3] + escape.extra;
            i = i + 1;
        }
        return acc;
    }

    private static int sharedOps(int[] array, Cell cell, int ops, int seed, int start, int stride) {
        int acc = seed;
        int i = 0;
        int mask = array.length - 1;
        while (i < ops) {
            int index = (start + i * stride) & mask;
            int value = array[index];
            value = value * 1103515245 + 12345;
            array[index] = value;
            cell.value = cell.value + value;
            acc = acc + value + cell.value;
            i = i + 1;
        }
        return acc;
    }

    static int sharedSingle(int[] array, Cell cell, int ops, int seed) {
        return sharedOps(array, cell, ops, seed, 0, 1);
    }

    static int sharedThreaded(int[] array, Cell[] cells, int totalOps, int threads, int seed)
            throws InterruptedException {
        SharedWorker[] workers = new SharedWorker[threads];
        int per = totalOps / threads;
        int t = 0;
        while (t < threads) {
            workers[t] = new SharedWorker(array, cells[t], per, seed + t, t, threads);
            t = t + 1;
        }
        t = 0;
        while (t < threads) {
            workers[t].start();
            t = t + 1;
        }
        t = 0;
        while (t < threads) {
            workers[t].join();
            t = t + 1;
        }
        int acc = 0;
        int completed = 0;
        t = 0;
        while (t < threads) {
            acc = acc + workers[t].acc;
            completed = completed + workers[t].completed;
            t = t + 1;
        }
        if (completed != totalOps) {
            throw new IllegalStateException("shared kernel incomplete");
        }
        return acc * 31 + completed;
    }

    static int monitorKernel(Counter counter, int totalOps, int threads) throws InterruptedException {
        MonitorWorker[] workers = new MonitorWorker[threads];
        int per = totalOps / threads;
        int t = 0;
        while (t < threads) {
            workers[t] = new MonitorWorker(counter, per);
            t = t + 1;
        }
        t = 0;
        while (t < threads) {
            workers[t].start();
            t = t + 1;
        }
        t = 0;
        while (t < threads) {
            workers[t].join();
            t = t + 1;
        }
        int completed = 0;
        t = 0;
        while (t < threads) {
            completed = completed + workers[t].completed;
            t = t + 1;
        }
        if (completed != totalOps || counter.value() != totalOps) {
            throw new IllegalStateException("monitor kernel incomplete");
        }
        return counter.value();
    }

    static final class SharedWorker extends Thread {
        private final int[] array;
        private final Cell cell;
        private final int ops;
        private final int seed;
        private final int start;
        private final int stride;
        int acc;
        int completed;

        SharedWorker(int[] array, Cell cell, int ops, int seed, int start, int stride) {
            this.array = array;
            this.cell = cell;
            this.ops = ops;
            this.seed = seed;
            this.start = start;
            this.stride = stride;
        }

        public void run() {
            acc = sharedOps(array, cell, ops, seed, start, stride);
            completed = ops;
        }
    }

    static final class MonitorWorker extends Thread {
        private final Counter counter;
        private final int ops;
        int completed;

        MonitorWorker(Counter counter, int ops) {
            this.counter = counter;
            this.ops = ops;
        }

        public void run() {
            int i = 0;
            while (i < ops) {
                counter.bump();
                i = i + 1;
            }
            completed = ops;
        }
    }

    static final class Counter {
        private int value;

        synchronized void bump() {
            value = value + 1;
        }

        synchronized int value() {
            return value;
        }
    }
}
