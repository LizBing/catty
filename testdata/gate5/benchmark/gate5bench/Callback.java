package gate5bench;

/**
 * The callback side of the mixed-engine crossing path. In the mixed
 * configuration this method stays AOT while {@link Hot#mix} is forced to the
 * interpreter by the build policy, so the AOT caller -> interpreted hot method
 * -> AOT callback chain crosses engines twice per operation.
 */
final class Callback {
    private Callback() {
    }

    public static int mix(int a, int b) {
        int x = a ^ (b + 0x9e3779b9);
        x = x * 1664525 + 1013904223;
        return x ^ (x >>> 16);
    }
}
