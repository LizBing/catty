package app;

public class Probe2 extends Base {
    public static Probe2 create() { return new Probe2(); }

    public int call() { return secret(); }

    @Override
    int tag() { return 2; }
}
