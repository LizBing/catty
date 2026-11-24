/** First explicit exit wins; a hook that calls exit can stall graceful shutdown. */
public final class ExitContention {
    static final Object LOCK = new Object();
    static volatile boolean hookStarted;
    static volatile boolean release;
    static volatile boolean exitingHookPrinted;

    static final class Holder implements Runnable {
        public void run() {
            System.out.println("hook-start");
            synchronized (LOCK) {
                hookStarted = true;
                LOCK.notifyAll();
                while (!release) {
                    try {
                        LOCK.wait();
                    } catch (InterruptedException ignored) {
                    }
                }
            }
            System.out.println("hook-end");
        }
    }

    static final class Second implements Runnable {
        public void run() {
            while (!hookStarted) {
                Thread.yield();
            }
            System.out.println("second");
            synchronized (LOCK) {
                release = true;
                LOCK.notifyAll();
            }
            System.exit(7);
            System.out.println("second-after");
        }
    }

    static final class ExitingHook implements Runnable {
        public void run() {
            System.out.println("pre-exit");
            exitingHookPrinted = true;
            System.exit(5);
            System.out.println("post-exit");
        }
    }

    static final class HalterHook implements Runnable {
        public void run() {
            while (!exitingHookPrinted) {
                Thread.yield();
            }
            try {
                Thread.sleep(40);
            } catch (InterruptedException ignored) {
            }
            System.out.println("halt-hook");
            Runtime.getRuntime().halt(21);
        }
    }

    public static void main(String[] args) {
        String mode = args[0];
        Runtime runtime = Runtime.getRuntime();
        if (mode.equals("status")) {
            runtime.addShutdownHook(new Thread(new Holder(), "holder"));
            new Thread(new Second(), "second-exit").start();
            System.out.println("main-exit");
            System.exit(3);
        } else if (mode.equals("hook-exit")) {
            runtime.addShutdownHook(new Thread(new ExitingHook(), "exiting-hook"));
            System.out.println("main-exit");
            System.exit(0);
        } else {
            runtime.addShutdownHook(new Thread(new ExitingHook(), "exiting-hook"));
            runtime.addShutdownHook(new Thread(new HalterHook(), "halter-hook"));
            System.out.println("main-exit");
            System.exit(0);
        }
        System.out.println("main-after");
    }
}
