/** Halt in the running state and halt during graceful shutdown. */
public final class HaltCase {
    static volatile boolean hookStarted;

    static final class Hook implements Runnable {
        public void run() {
            System.out.println("hook-start");
            hookStarted = true;
            for (;;) {
                Thread.yield();
            }
        }
    }

    static final class Blocked implements Runnable {
        public void run() {
            for (;;) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException ignored) {
                }
            }
        }
    }

    static final class Halter implements Runnable {
        public void run() {
            while (!hookStarted) {
                Thread.yield();
            }
            System.out.println("halter");
            Runtime.getRuntime().halt(23);
        }
    }

    public static void main(String[] args) {
        String mode = args[0];
        Runtime runtime = Runtime.getRuntime();
        runtime.addShutdownHook(new Thread(new Hook(), "hook"));
        if (mode.equals("before")) {
            new Thread(new Blocked(), "blocked").start();
            System.out.println("before");
            runtime.halt(19);
        } else {
            new Thread(new Halter(), "halter").start();
            System.out.println("before");
            System.exit(3);
        }
        System.out.println("after");
    }
}
