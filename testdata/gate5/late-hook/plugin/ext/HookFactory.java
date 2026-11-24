package ext;

/** Plugin-supplied hook factory loaded only through a late admission. */
public final class HookFactory {
    private HookFactory() {}

    public static Thread make() {
        return new PluginHook();
    }
}
