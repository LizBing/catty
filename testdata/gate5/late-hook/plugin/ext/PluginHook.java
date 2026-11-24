package ext;

/** Plugin-supplied hook body: interpreted, own logical thread, identity checked. */
public final class PluginHook extends Thread {
    public PluginHook() {
        super("plugin-hook");
    }

    public void run() {
        System.out.println(this == Thread.currentThread() ? "plugin-identity=1" : "plugin-identity=0");
        System.out.println("plugin-hook");
    }
}
