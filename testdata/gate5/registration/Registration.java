/** Hook registry rules, identity, retired peers and frozen-registry probes. */
public final class Registration {
    static final Object LOCK = new Object();
    static volatile int equalRuns;
    static volatile boolean liveRelease;
    static volatile boolean liveStarted;

    static final class Blocking implements Runnable {
        public void run() {
            liveStarted = true;
            while (!liveRelease) {
                Thread.yield();
            }
        }
    }

    static final class RemovedHook extends Thread {
        public void run() {
            System.out.println("removed-ran");
        }
    }

    static final class RetiredHook extends Thread {
        public void run() {
            System.out.println("retired-ran");
        }
    }

    static final class EqualHook extends Thread {
        public void run() {
            synchronized (LOCK) {
                equalRuns++;
                LOCK.notifyAll();
            }
        }

        public boolean equals(Object other) {
            return other instanceof EqualHook;
        }
    }

    static final class Aggregator extends Thread {
        Aggregator() {
            super("aggregator");
        }

        public void run() {
            synchronized (LOCK) {
                while (equalRuns < 2) {
                    try {
                        LOCK.wait();
                    } catch (InterruptedException ignored) {
                    }
                }
            }
            System.out.println("equal-hooks=2");
            Runtime runtime = Runtime.getRuntime();
            try {
                runtime.addShutdownHook(null);
                System.out.println("frozen-add-null-accepted");
            } catch (NullPointerException wrong) {
                System.out.println("frozen-add-null-npe");
            } catch (IllegalStateException expected) {
                System.out.println("frozen-add-null-ise");
            }
            try {
                runtime.addShutdownHook(new Thread());
                System.out.println("frozen-add-thread-accepted");
            } catch (IllegalStateException expected) {
                System.out.println("frozen-add-thread-ise");
            }
            try {
                runtime.removeShutdownHook(null);
                System.out.println("frozen-remove-null-accepted");
            } catch (NullPointerException wrong) {
                System.out.println("frozen-remove-null-npe");
            } catch (IllegalStateException expected) {
                System.out.println("frozen-remove-null-ise");
            }
        }
    }

    public static void main(String[] args) throws InterruptedException {
        Runtime runtime = Runtime.getRuntime();
        System.out.println(runtime == Runtime.getRuntime() ? "singleton=1" : "singleton=0");
        System.out.println(runtime == null ? "nonnull=0" : "nonnull=1");

        RemovedHook removed = new RemovedHook();
        runtime.addShutdownHook(removed);
        System.out.println("add=1");
        System.out.println(runtime.removeShutdownHook(removed) ? "remove=1" : "remove=0");
        System.out.println(runtime.removeShutdownHook(removed) ? "remove-absent=1" : "remove-absent=0");

        try {
            runtime.addShutdownHook(null);
            System.out.println("add-null-accepted");
        } catch (NullPointerException expected) {
            System.out.println("add-null-npe");
        }
        try {
            runtime.removeShutdownHook(null);
            System.out.println("remove-null-accepted");
        } catch (NullPointerException expected) {
            System.out.println("remove-null-npe");
        }

        EqualHook duplicate = new EqualHook();
        runtime.addShutdownHook(duplicate);
        try {
            runtime.addShutdownHook(duplicate);
            System.out.println("duplicate-accepted");
        } catch (IllegalArgumentException expected) {
            System.out.println("duplicate-iae");
        }

        Thread live = new Thread(new Blocking(), "live");
        live.start();
        while (!liveStarted) {
            Thread.yield();
        }
        try {
            runtime.addShutdownHook(live);
            System.out.println("live-accepted");
        } catch (IllegalArgumentException expected) {
            System.out.println("live-iae");
        }
        liveRelease = true;
        live.join();

        RetiredHook retired = new RetiredHook();
        retired.start();
        retired.join();
        runtime.addShutdownHook(retired);
        System.out.println("retired-registered");

        EqualHook left = new EqualHook();
        EqualHook right = new EqualHook();
        runtime.addShutdownHook(left);
        runtime.addShutdownHook(right);
        System.out.println("equal-registered");

        runtime.addShutdownHook(new Aggregator());
        System.exit(0);
    }
}
