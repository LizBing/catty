public class WaitSet {
    private static final Object LOCK = new Object();
    private static volatile boolean sleeping;
    private static boolean sleepInterrupted;

    static class Sleeper extends Thread {
        public void run() {
            synchronized (LOCK) {
                sleeping = true;
                try {
                    Thread.sleep(9223372036854775807L, 999999);
                    throw new RuntimeException("long sleep returned without interrupt");
                } catch (InterruptedException expected) {
                    sleepInterrupted = Thread.holdsLock(LOCK) && !Thread.interrupted();
                }
            }
        }
    }

    static class Trial {
        private final Object monitor = new Object();
        private volatile int ready;
        private boolean signalled;
        private boolean normalWake;

        class Waiter extends Thread {
            boolean interrupted;
            boolean pending;

            public void run() {
                synchronized (monitor) {
                    synchronized (monitor) {
                        ready++;
                        try {
                            while (!signalled) {
                                monitor.wait();
                            }
                            pending = Thread.currentThread().isInterrupted();
                            normalWake = true;
                            monitor.notifyAll();
                        } catch (InterruptedException expected) {
                            interrupted = true;
                            Check.that(!Thread.currentThread().isInterrupted(), "wait failed to clear interrupt");
                        }
                        Check.that(Thread.holdsLock(monitor), "wait returned without its monitor");
                    }
                    Check.that(Thread.holdsLock(monitor), "wait failed to restore recursion depth");
                }
            }
        }

        private void run(boolean interruptFirst) throws InterruptedException {
            Waiter a = new Waiter();
            Waiter b = new Waiter();
            a.start();
            b.start();
            while (ready != 2) {
                Thread.yield();
            }
            synchronized (monitor) {
                signalled = true;
                if (interruptFirst) {
                    a.interrupt();
                    monitor.notify();
                } else {
                    monitor.notify();
                    a.interrupt();
                }
                while (!normalWake) {
                    monitor.wait();
                }
                monitor.notifyAll();
            }
            a.join();
            b.join();
            Check.that(a.interrupted != a.pending && !b.interrupted && !b.pending,
                    "notify/interrupt returned a forbidden outcome");
        }
    }

    private static void ownershipAndTimeouts() throws InterruptedException {
        try {
            LOCK.notify();
            throw new RuntimeException("notify without ownership");
        } catch (IllegalMonitorStateException expected) {
            // Ownership checks apply to all three operations.
        }
        try {
            LOCK.notifyAll();
            throw new RuntimeException("notifyAll without ownership");
        } catch (IllegalMonitorStateException expected) {
        }
        try {
            LOCK.wait(1);
            throw new RuntimeException("wait without ownership");
        } catch (IllegalMonitorStateException expected) {
        }
        synchronized (LOCK) {
            try {
                LOCK.wait(0, 1000000);
                throw new RuntimeException("invalid nanos accepted");
            } catch (IllegalArgumentException expected) {
            }
            Thread.currentThread().interrupt();
            try {
                LOCK.wait();
                throw new RuntimeException("pre-interrupted wait returned normally");
            } catch (InterruptedException expected) {
                Check.that(Thread.holdsLock(LOCK) && !Thread.interrupted(), "pre-interrupted wait semantics");
            }
            LOCK.wait(0, 1);
            Check.that(Thread.holdsLock(LOCK), "timed wait lost monitor");
            Thread.currentThread().interrupt();
            try {
                Thread.sleep(10000);
                throw new RuntimeException("pre-interrupted sleep returned normally");
            } catch (InterruptedException expected) {
                Check.that(Thread.holdsLock(LOCK) && !Thread.interrupted(), "sleep interrupt semantics");
            }
        }
    }

    public static void main(String[] args) throws InterruptedException {
        ownershipAndTimeouts();
        Sleeper sleeper = new Sleeper();
        sleeper.start();
        while (!sleeping) {
            Thread.yield();
        }
        sleeper.interrupt();
        sleeper.join();
        Check.that(sleepInterrupted, "sleep did not retain ownership and clear the interrupt");
        for (int i = 0; i < 16; i++) {
            new Trial().run((i & 1) == 0);
        }
        System.out.println("wait-set");
    }
}
