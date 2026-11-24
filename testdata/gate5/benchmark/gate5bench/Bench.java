package gate5bench;

import java.lang.reflect.Method;

/**
 * Command implementations. Every command prints the deterministic record
 * protocol documented in {@link Report}; the Python runner in
 * {@code tests/gate5_bench.py} parses it and checks it against the OpenJDK 21
 * oracle.
 */
final class Bench {
    private Bench() {
    }

    private static void check(String name, int value) {
        Report.word("CHECK");
        Report.word(name);
        Report.number(value);
    }

    private static void batch(long nanos) {
        Report.word("NS");
        Report.bignumber(nanos);
    }

    private static void kernelHeader(String name, int ops, int batches, int warmup) {
        Report.word("KERNEL");
        Report.word(name);
        Report.number(ops);
        Report.number(batches);
        Report.number(warmup);
    }

    private static Calls.Base[] receivers() {
        Calls.Base[] receivers = new Calls.Base[2];
        receivers[0] = new Calls.Alpha();
        receivers[1] = new Calls.Beta();
        return receivers;
    }

    private static Cell[] cells(int count) {
        Cell[] cells = new Cell[count];
        int i = 0;
        while (i < count) {
            cells[i] = new Cell();
            i = i + 1;
        }
        return cells;
    }

    private static int[] slots() {
        return new int[Sizes.SHARED_SLOTS];
    }

    /** Correctness command: every kernel once, deterministic checksums only. */
    static void verify() throws Exception {
        Report.begin();
        Report.word("COMMAND");
        Report.word("verify");

        Node graph = Workload.buildGraph(Sizes.APP_NODES, Sizes.SEED);
        check("app", Workload.process(graph, Sizes.APP_NODES, Sizes.APP_ROUNDS));

        check("calls-static", Calls.staticKernel(Sizes.CALL_STATIC_OPS, Sizes.SEED));
        check("calls-virtual", Calls.virtualKernel(receivers(), Sizes.CALL_VIRTUAL_OPS, Sizes.SEED));
        check("calls-mixed", Calls.crossingKernel(Sizes.CALL_CROSSING_OPS, Sizes.SEED));

        check("reflect-lookup", Reflection.lookupKernel(Sizes.REFLECT_LOOKUP_OPS, Sizes.SEED));
        Method method = Reflection.lookupCombine();
        Object[] args = new Object[]{Integer.valueOf(0), Integer.valueOf(0)};
        check("reflect-invoke", Reflection.invokeKernel(method, args, Sizes.REFLECT_INVOKE_OPS, Sizes.SEED));

        Escape[] ring = new Escape[Sizes.ALLOC_OPS];
        check("memory-alloc", Memory.allocKernel(ring, Sizes.ALLOC_OPS, Sizes.SEED));

        Cell[] cells = cells(1);
        check("shared-1", Memory.sharedSingle(slots(), cells[0], Sizes.SHARED_OPS, Sizes.SEED));
        check("shared-4", Memory.sharedThreaded(slots(), cells(Sizes.SHARED_THREADS),
                Sizes.SHARED_OPS, Sizes.SHARED_THREADS, Sizes.SEED));

        int monitor = Memory.monitorKernel(new Memory.Counter(), Sizes.MONITOR_OPS, Sizes.MONITOR_THREADS);
        if (monitor != Sizes.MONITOR_OPS) {
            throw new IllegalStateException("monitor work mismatch");
        }
        check("monitor", monitor * 31 + monitor);
        Report.word("WORK");
        Report.number(monitor);

        Report.end();
    }

    /**
     * Cold command: one app workload pass, then a single verified result record.
     * The runner timestamps receipt of the checksum line from process launch.
     */
    static void cold() {
        Node graph = Workload.buildGraph(Sizes.APP_NODES, Sizes.SEED);
        int result = Workload.process(graph, Sizes.APP_NODES, Sizes.APP_ROUNDS);
        Report.begin();
        Report.word("COMMAND");
        Report.word("cold");
        Report.word("RESULT");
        Report.number(result);
        Report.number(Sizes.APP_NODES * Sizes.APP_ROUNDS);
        Report.word("app");
        Report.end();
    }

    /** Steady command: fixed warmup then timed batches for every kernel. */
    static void steady(int batches, int warmup) throws Exception {
        Report.begin();
        Report.word("COMMAND");
        Report.word("steady");
        Report.number(10);

        steadyApp(batches, warmup);
        steadyCalls(batches, warmup);
        steadyReflection(batches, warmup);
        steadyMemory(batches, warmup);

        Report.end();
    }

    private static void steadyApp(int batches, int warmup) {
        int ops = Sizes.APP_NODES * Sizes.APP_ROUNDS;
        kernelHeader("app", ops, batches, warmup);
        Node graph = Workload.buildGraph(Sizes.APP_NODES, Sizes.SEED);
        int state = Sizes.SEED;
        int w = 0;
        while (w < warmup) {
            state = Workload.process(graph, Sizes.APP_NODES, Sizes.APP_ROUNDS);
            w = w + 1;
        }
        int checksum = 0;
        int b = 0;
        while (b < batches) {
            long t0 = System.nanoTime();
            int result = Workload.process(graph, Sizes.APP_NODES, Sizes.APP_ROUNDS);
            long t1 = System.nanoTime();
            state = result;
            checksum = checksum + result;
            batch(t1 - t0);
            b = b + 1;
        }
        check("app", checksum + (state & 1));
    }

    private static void steadyCalls(int batches, int warmup) {
        {
            int ops = Sizes.CALL_STATIC_OPS;
            kernelHeader("calls-static", ops, batches, warmup);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Calls.staticKernel(ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Calls.staticKernel(ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("calls-static", checksum);
        }
        {
            int ops = Sizes.CALL_VIRTUAL_OPS;
            kernelHeader("calls-virtual", ops, batches, warmup);
            Calls.Base[] receivers = receivers();
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Calls.virtualKernel(receivers, ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Calls.virtualKernel(receivers, ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("calls-virtual", checksum);
        }
        {
            int ops = Sizes.CALL_CROSSING_OPS;
            kernelHeader("calls-mixed", ops, batches, warmup);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Calls.crossingKernel(ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Calls.crossingKernel(ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("calls-mixed", checksum);
        }
    }

    private static void steadyReflection(int batches, int warmup) throws Exception {
        {
            int ops = Sizes.REFLECT_LOOKUP_OPS;
            kernelHeader("reflect-lookup", ops, batches, warmup);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Reflection.lookupKernel(ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Reflection.lookupKernel(ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("reflect-lookup", checksum);
        }
        {
            int ops = Sizes.REFLECT_INVOKE_OPS;
            kernelHeader("reflect-invoke", ops, batches, warmup);
            Method method = Reflection.lookupCombine();
            Object[] args = new Object[]{Integer.valueOf(0), Integer.valueOf(0)};
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Reflection.invokeKernel(method, args, ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Reflection.invokeKernel(method, args, ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("reflect-invoke", checksum);
        }
    }

    private static void steadyMemory(int batches, int warmup) throws Exception {
        {
            int ops = Sizes.ALLOC_OPS;
            kernelHeader("memory-alloc", ops, batches, warmup);
            Escape[] ring = new Escape[Sizes.ALLOC_RING];
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Memory.allocKernel(ring, ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Memory.allocKernel(ring, ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("memory-alloc", checksum);
        }
        {
            int ops = Sizes.SHARED_OPS;
            kernelHeader("shared-1", ops, batches, warmup);
            int[] array = slots();
            Cell[] cells = cells(1);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Memory.sharedSingle(array, cells[0], ops, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Memory.sharedSingle(array, cells[0], ops, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("shared-1", checksum);
        }
        {
            int ops = Sizes.SHARED_OPS;
            kernelHeader("shared-4", ops, batches, warmup);
            int[] array = slots();
            Cell[] cells = cells(Sizes.SHARED_THREADS);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                state = Memory.sharedThreaded(array, cells, ops, Sizes.SHARED_THREADS, state);
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                long t0 = System.nanoTime();
                int result = Memory.sharedThreaded(array, cells, ops, Sizes.SHARED_THREADS, state);
                long t1 = System.nanoTime();
                state = result;
                checksum = checksum + result;
                batch(t1 - t0);
                b = b + 1;
            }
            check("shared-4", checksum);
        }
        {
            int ops = Sizes.MONITOR_OPS;
            kernelHeader("monitor", ops, batches, warmup);
            int state = Sizes.SEED;
            int w = 0;
            while (w < warmup) {
                Memory.Counter warm = new Memory.Counter();
                int done = Memory.monitorKernel(warm, ops, Sizes.MONITOR_THREADS);
                state = state + done;
                w = w + 1;
            }
            int checksum = 0;
            int b = 0;
            while (b < batches) {
                Memory.Counter counter = new Memory.Counter();
                long t0 = System.nanoTime();
                int done = Memory.monitorKernel(counter, ops, Sizes.MONITOR_THREADS);
                long t1 = System.nanoTime();
                state = state + done;
                checksum = checksum + done * 31 + counter.value();
                batch(t1 - t0);
                b = b + 1;
            }
            check("monitor", checksum + (state & 255));
        }
    }

    /** Peak-RSS sample command: the memory kernels with retained allocations. */
    static void rss() throws Exception {
        Report.begin();
        Report.word("COMMAND");
        Report.word("rss");

        Escape[] ring = new Escape[Sizes.ALLOC_RSS_OPS];
        int alloc = Memory.allocKernel(ring, Sizes.ALLOC_RSS_OPS, Sizes.SEED);
        check("memory-alloc", alloc);

        Cell[] single = cells(1);
        check("shared-1", Memory.sharedSingle(slots(), single[0], Sizes.SHARED_OPS, Sizes.SEED));
        check("shared-4", Memory.sharedThreaded(slots(), cells(Sizes.SHARED_THREADS),
                Sizes.SHARED_OPS, Sizes.SHARED_THREADS, Sizes.SEED));

        Memory.Counter counter = new Memory.Counter();
        int monitor = Memory.monitorKernel(counter, Sizes.MONITOR_OPS, Sizes.MONITOR_THREADS);
        check("monitor", monitor * 31 + monitor);
        Report.word("WORK");
        Report.number(monitor);

        int live = 0;
        int i = 0;
        while (i < ring.length) {
            live = live + ring[i].payload[3];
            i = i + 1;
        }
        check("ring-live", live);

        Report.end();
    }
}
