final class Check {
    static void that(boolean condition, String message) {
        if (!condition) {
            throw new RuntimeException(message);
        }
    }
}
