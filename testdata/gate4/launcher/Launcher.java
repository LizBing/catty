public class Launcher {
    private static final Object GATE = new Object();
    private static boolean daemonReady;
    private static boolean release;
    private static boolean mainWillThrow;
    private static Thread mainThread;

    static class Daemon extends Thread {
        public void run() {
            synchronized (GATE) {
                Check.that(new Thread().isDaemon(), "daemon status must be inherited");
                daemonReady = true;
                GATE.notifyAll();
                for (;;) {
                    try {
                        GATE.wait();
                    } catch (InterruptedException e) {
                        throw new RuntimeException("daemon interrupted by launcher");
                    }
                }
            }
        }
    }

    static class Descendant extends Thread {
        private final Thread parent;

        Descendant(Thread parent) {
            this.parent = parent;
        }

        public void run() {
            while (parent.isAlive()) {
                Thread.yield();
            }
            System.out.println("descendant");
        }
    }

    static class Worker extends Thread {
        public void run() {
            synchronized (GATE) {
                while (!release) {
                    try {
                        GATE.wait();
                    } catch (InterruptedException e) {
                        throw new RuntimeException("worker interrupted by launcher");
                    }
                }
            }
            while (mainThread.isAlive()) {
                Thread.yield();
            }
            System.out.println("worker");
            new Descendant(this).start();
            if (!mainWillThrow) {
                throw new RuntimeException("worker failed");
            }
        }
    }

    public static void main(String[] args) throws InterruptedException {
        mainThread = Thread.currentThread();
        mainWillThrow = args.length > 0;
        Daemon daemon = new Daemon();
        daemon.setDaemon(true);
        daemon.start();
        try {
            daemon.setDaemon(false);
            throw new RuntimeException("changed daemon status of live thread");
        } catch (IllegalThreadStateException expected) {
        }
        synchronized (GATE) {
            while (!daemonReady) {
                GATE.wait();
            }
            new Worker().start();
            System.out.println("main");
            release = true;
            GATE.notifyAll();
        }
        if (mainWillThrow) {
            throw new RuntimeException("main failed");
        }
    }
}
