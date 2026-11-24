public class Lifecycle {
    private static final Object GATE = new Object();
    private static final Object JOIN_GATE = new Object();
    private static Thread main;
    private static volatile boolean started;
    private static volatile boolean joining;
    private static boolean release;
    private static Thread observed;
    private static int input;
    private static int result;
    private static int directCalls;
    private static boolean joinInterrupted;
    private static volatile boolean startRace;

    static class Worker extends Thread {
        public void run() {
            observed = Thread.currentThread();
            if (observed == main) {
                directCalls++;
                return;
            }
            synchronized (GATE) {
                started = true;
                while (!release) {
                    try {
                        GATE.wait();
                    } catch (InterruptedException e) {
                        throw new RuntimeException("unexpected worker interrupt");
                    }
                }
            }
            result = input + 1;
        }
    }

    static class Joiner implements Runnable {
        private final Thread target;

        Joiner(Thread target) {
            this.target = target;
        }

        public void run() {
            synchronized (JOIN_GATE) {
                joining = true;
                try {
                    target.join();
                } catch (InterruptedException e) {
                    joinInterrupted = !Thread.currentThread().isInterrupted()
                            && Thread.holdsLock(JOIN_GATE);
                }
            }
        }
    }

    static class Counted extends Thread {
        private int runs;

        Counted() {
            super("counted");
        }

        public synchronized void run() {
            runs++;
        }
    }

    static class Starter extends Thread {
        private final Counted target;
        private boolean won;

        Starter(Counted target) {
            this.target = target;
        }

        public void run() {
            while (!startRace) {
                Thread.yield();
            }
            try {
                target.start();
                won = true;
            } catch (IllegalThreadStateException expected) {
            }
        }
    }

    private static void simultaneousStart() throws InterruptedException {
        Counted target = new Counted();
        long id = target.threadId();
        Starter left = new Starter(target);
        Starter right = new Starter(target);
        left.start();
        right.start();
        startRace = true;
        left.join();
        right.join();
        target.join();
        Check.that(left.won != right.won && target.runs == 1, "two first starts succeeded");
        Check.that(id > 0 && id == target.threadId() && id != main.threadId(), "incorrect stable thread ID");
    }

    public static void main(String[] args) throws InterruptedException {
        main = Thread.currentThread();
        Worker worker = new Worker();
        Check.that(main == Thread.currentThread() && main.isAlive(), "main thread identity");
        worker.join();
        Check.that(!worker.isAlive(), "joining a new thread");
        worker.run();
        Check.that(directCalls == 1 && !worker.isAlive(), "direct run must remain synchronous");
        input = 41;
        worker.start();
        while (!started) {
            Thread.yield();
        }
        Check.that(worker.isAlive(), "worker should be alive at gate");
        worker.join(1, 1);
        Check.that(worker.isAlive(), "timed join should return with worker still alive");
        Thread joiner = new Thread(new Joiner(worker));
        joiner.start();
        while (!joining) {
            Thread.yield();
        }
        joiner.interrupt();
        joiner.join();
        Check.that(joinInterrupted, "join interrupt should clear status and retain other monitors");

        Thread contender = new Thread(new Runnable() {
            public void run() {
                try {
                    worker.start();
                } catch (IllegalThreadStateException e) {
                    return;
                }
                throw new RuntimeException("a live thread started twice");
            }
        });
        contender.start();
        contender.join();
        synchronized (GATE) {
            release = true;
            GATE.notifyAll();
        }
        worker.join();
        Check.that(observed == worker && result == 42 && !worker.isAlive(), "worker identity and join publication");
        try {
            worker.start();
            throw new RuntimeException("terminated thread restarted");
        } catch (IllegalThreadStateException expected) {
            Thread.currentThread().interrupt();
            Check.that(main.isInterrupted() && main.isInterrupted(), "status query cleared interrupt");
            Check.that(Thread.interrupted() && !Thread.interrupted(), "test-and-clear interrupt");
        }
        Thread inactive = new Thread();
        Check.that(!inactive.isDaemon(), "main is non-daemon");
        inactive.setDaemon(true);
        Check.that(inactive.isDaemon(), "new thread daemon setting");
        simultaneousStart();
        try {
            worker.join(-1);
            throw new RuntimeException("negative join accepted");
        } catch (IllegalArgumentException expected) {
            System.out.println("lifecycle");
        }
    }
}
