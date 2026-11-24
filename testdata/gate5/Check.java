/** Shared Gate 5 assertion helper; throws a named Java error on mismatch. */
public final class Check {
    private Check() {}

    public static void that(boolean condition, String message) {
        if (!condition) {
            throw new RuntimeException(message);
        }
    }
}
