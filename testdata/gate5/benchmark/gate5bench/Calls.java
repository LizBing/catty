package gate5bench;

/**
 * Static/direct and virtual dispatch kernels, plus the mixed-policy crossing
 * kernel. Every kernel has a fixed operation count and a deterministic
 * checksum.
 */
final class Calls {
    private Calls() {
    }

    /** The direct callee; the caller uses {@code invokestatic}. */
    static int directTarget(int acc, int i) {
        int x = acc + (i * 3 + 1);
        x = x * 1103515245 + 12345;
        return x ^ (x >>> 13);
    }

    static int staticKernel(int ops, int seed) {
        int acc = seed;
        int i = 0;
        while (i < ops) {
            acc = directTarget(acc, i);
            i = i + 1;
        }
        return acc;
    }

    /** Two concrete receiver types for one virtual call site. */
    static abstract class Base {
        abstract int pick(int acc, int i);
    }

    static final class Alpha extends Base {
        int pick(int acc, int i) {
            int x = acc + i * 5 + 11;
            return x ^ (x >>> 15);
        }
    }

    static final class Beta extends Base {
        int pick(int acc, int i) {
            int x = acc ^ (i * 9 + 29);
            x = x * 1664525 + 1013904223;
            return x + (x >>> 7);
        }
    }

    static int virtualKernel(Base[] receivers, int ops, int seed) {
        int acc = seed;
        int i = 0;
        int n = receivers.length;
        while (i < ops) {
            Base receiver = receivers[i & (n - 1)];
            acc = receiver.pick(acc, i);
            i = i + 1;
        }
        return acc;
    }

    /**
     * AOT caller -> interpreted {@link Hot#mix} -> AOT {@link Callback}: the
     * mixed-engine crossing kernel.
     */
    static int crossingKernel(int ops, int seed) {
        int acc = seed;
        int i = 0;
        while (i < ops) {
            acc = Hot.mix(acc, i);
            i = i + 1;
        }
        return acc;
    }
}
