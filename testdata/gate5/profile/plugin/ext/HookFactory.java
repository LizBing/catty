package ext;

/** Otherwise valid late class that references a Gate 5-only bootstrap member. */
public final class HookFactory {
    private HookFactory() {}

    public static Thread make() {
        System.out.println(System.nanoTime());
        return new Thread();
    }
}
