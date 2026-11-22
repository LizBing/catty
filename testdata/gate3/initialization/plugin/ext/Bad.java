package ext;

public class Bad {
    static {
        if (app.Probe.supplied != null) {
            throw (RuntimeException) app.Probe.supplied;
        }
    }

    public static int value() { return 9; }
}
