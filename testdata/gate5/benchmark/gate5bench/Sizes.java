package gate5bench;

/**
 * Fixed workload sizes shared by every configuration. No value is derived from
 * the wall clock, the environment, or the number of processors, so the same
 * class files perform the same externally visible work on OpenJDK 21 and on
 * each Catty engine configuration.
 */
final class Sizes {
    private Sizes() {
    }

    static final int SEED = 12345;

    // Application object graph.
    static final int APP_NODES = 1024;
    static final int APP_ROUNDS = 1;

    // Calls.
    static final int CALL_STATIC_OPS = 20000;
    static final int CALL_VIRTUAL_OPS = 20000;
    static final int CALL_CROSSING_OPS = 4000;

    // Reflection.
    static final int REFLECT_LOOKUP_OPS = 300;
    static final int REFLECT_INVOKE_OPS = 2000;

    // Memory and threads.
    static final int ALLOC_OPS = 4096;
    static final int ALLOC_RING = 1024;
    static final int ALLOC_RSS_OPS = 32768;
    static final int SHARED_OPS = 8192;
    static final int SHARED_THREADS = 4;
    static final int SHARED_SLOTS = 512;
    static final int MONITOR_OPS = 4096;
    static final int MONITOR_THREADS = 4;

    // Retained graph for the explicit-GC observation.
    static final int RETAINED_NODES = 12000;

    // Warmup and measured batch counts. The full command is used by the
    // evidence run; the smoke command only by the runner smoke check.
    static final int FULL_BATCHES = 30;
    static final int FULL_WARMUP = 10;
    static final int SMOKE_BATCHES = 5;
    static final int SMOKE_WARMUP = 2;
}
