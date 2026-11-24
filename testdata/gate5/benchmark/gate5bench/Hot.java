package gate5bench;

/**
 * The designated hot method of the mixed configuration. The build policy in
 * {@code tests/gate5_bench.py} forces exactly this method to interpretation
 * while its AOT callers and the AOT {@link Callback} stay compiled.
 */
final class Hot {
    private Hot() {
    }

    public static int mix(int a, int b) {
        int x = Callback.mix(a, b);
        return x + (b & 31);
    }
}
