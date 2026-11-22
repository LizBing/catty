package ext;

public class NestOuter {
    private static int value() { return 5; }

    public static int read() { return NestInner.access(); }

    static class NestInner {
        static int access() { return value(); }
    }
}
