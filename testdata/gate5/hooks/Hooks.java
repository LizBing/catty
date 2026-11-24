/** Concurrent hooks, a barrier, a daemon hook, virtual run overrides and a throwing hook. */
public final class Hooks {
    static final Object BARRIER = new Object();
    static int arrived;

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

    static final class HookA extends Thread {
        HookA() {
            super("hook-a");
        }

        public void run() {
            System.out.println(this == Thread.currentThread() ? "identity-a=1" : "identity-a=0");
            synchronized (BARRIER) {
                arrived++;
                BARRIER.notifyAll();
                while (arrived < 2) {
                    try {
                        BARRIER.wait();
                    } catch (InterruptedException ignored) {
                    }
                }
            }
            new Thread(new Blocked(), "hook-worker").start();
            System.out.println("hook-a");
        }
    }

    static final class HookB extends Thread {
        HookB() {
            super("hook-b");
        }

        public void run() {
            System.out.println(this == Thread.currentThread() ? "identity-b=1" : "identity-b=0");
            synchronized (BARRIER) {
                arrived++;
                BARRIER.notifyAll();
                while (arrived < 2) {
                    try {
                        BARRIER.wait();
                    } catch (InterruptedException ignored) {
                    }
                }
            }
            System.out.println("hook-b");
        }
    }

    static final class Throwing extends Thread {
        Throwing() {
            super("bad-hook");
        }

        public void run() {
            System.out.println("throwing");
            throw new RuntimeException("hook-boom");
        }
    }

    public static void main(String[] args) {
        Runtime runtime = Runtime.getRuntime();
        runtime.addShutdownHook(new HookA());
        HookB hookB = new HookB();
        hookB.setDaemon(true);
        runtime.addShutdownHook(hookB);
        runtime.addShutdownHook(new Throwing());
        System.out.println("main");
        System.exit(7);
    }
}
