/** Exit from main-class initialization before main runs. */
public final class ClinitExit {
    static {
        Runtime runtime = Runtime.getRuntime();
        runtime.addShutdownHook(new Thread(new ExitCase.Hook(), "hook"));
        System.out.println("clinit");
        System.exit(7);
    }

    public static void main(String[] args) {
        System.out.println("main");
    }
}
