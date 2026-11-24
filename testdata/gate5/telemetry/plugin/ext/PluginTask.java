package ext;

/** Late-admitted interpreted worker used by the telemetry fixture. */
public final class PluginTask extends Thread {
    public PluginTask() {
        super("plugin-task");
    }

    public static Thread make() {
        return new PluginTask();
    }

    public void run() {
        System.out.println("plugin");
    }
}
