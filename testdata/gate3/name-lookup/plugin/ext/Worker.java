package ext;

public class Worker {
    static {
        app.Probe.mark();
    }

    public static int magic() { return 73; }
}
