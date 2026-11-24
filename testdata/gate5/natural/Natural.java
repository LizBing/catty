/** Natural termination: non-daemon drain, a descendant, a main joiner, daemon and hooks. */
public final class Natural {
    static Thread mainThread;
    static final Object START = new Object();
    static volatile boolean releaseWorker;
    static volatile boolean daemonReady;
    static volatile boolean descendantDone;
    static volatile boolean workerDone;
    static volatile boolean waitForever;
    static volatile boolean hookWaiting;
    static final Object HOOK_WAIT = new Object();

    static final class Hook implements Runnable {
        public void run() {
            if (waitForever) {
                System.out.println("hook-wait");
                hookWaiting = true;
                synchronized (HOOK_WAIT) {
                    for (;;) {
                        try {
                            HOOK_WAIT.wait();
                        } catch (InterruptedException ignored) {
                        }
                    }
                }
            }
            while (!workerDone) {
                Thread.yield();
            }
            System.out.println("hook");
        }
    }

    static final class Halter implements Runnable {
        public void run() {
            while (!hookWaiting) {
                Thread.yield();
            }
            System.out.println("halter");
            Runtime.getRuntime().halt(31);
        }
    }

    static final class Daemon implements Runnable {
        public void run() {
            synchronized (START) {
                daemonReady = true;
                START.notifyAll();
            }
            for (;;) {
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException ignored) {
                }
            }
        }
    }

    static final class Descendant implements Runnable {
        public void run() {
            System.out.println("descendant");
            descendantDone = true;
        }
    }

    static final class Joiner implements Runnable {
        public void run() {
            try {
                mainThread.join();
            } catch (InterruptedException ignored) {
            }
            while (!descendantDone) {
                Thread.yield();
            }
            System.out.println("joined");
        }
    }

    static final class Worker implements Runnable {
        public void run() {
            while (!releaseWorker) {
                Thread.yield();
            }
            Thread child = new Thread(new Descendant(), "descendant");
            child.start();
            try {
                child.join();
            } catch (InterruptedException ignored) {
            }
            Thread joiner = new Thread(new Joiner(), "joiner");
            joiner.start();
            try {
                joiner.join();
            } catch (InterruptedException ignored) {
            }
            System.out.println("worker");
            workerDone = true;
        }
    }

    public static void main(String[] args) throws InterruptedException {
        if (args.length > 0 && (args[0].equals("wait") || args[0].equals("wait-halt"))) {
            waitForever = true;
            Runtime.getRuntime().addShutdownHook(new Thread(new Hook(), "hook"));
            if (args[0].equals("wait-halt")) {
                Thread halter = new Thread(new Halter(), "halter");
                halter.setDaemon(true);
                halter.start();
            }
            System.out.println("main");
            return;
        }
        mainThread = Thread.currentThread();
        Runtime.getRuntime().addShutdownHook(new Thread(new Hook(), "hook"));
        Thread daemon = new Thread(new Daemon(), "daemon");
        daemon.setDaemon(true);
        daemon.start();
        synchronized (START) {
            while (!daemonReady) {
                START.wait();
            }
        }
        new Thread(new Worker(), "worker").start();
        System.out.println("main");
        releaseWorker = true;
        if (args.length > 0) {
            throw new RuntimeException("main failed");
        }
    }
}
