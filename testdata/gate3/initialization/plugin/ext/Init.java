package ext;

public class Init extends InitBase {
    static {
        app.Probe.mark();
    }

    public static int value() { return 5; }
}
