package app;

public class Probe {
    public static int count;

    public static Throwable supplied;

    public static void mark() { count++; }

    public static int get() { return count; }

    public static void reset() { count = 0; }
}
