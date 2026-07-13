public class InterruptWait {
    static final Object LOCK = new Object();
    static boolean waiting;

    static class Worker extends Thread {
        public void run() {
            synchronized (LOCK) {
                waiting = true;
                LOCK.notifyAll();
                try {
                    LOCK.wait();
                    System.out.println("missing");
                } catch (InterruptedException expected) {
                    System.out.println(Thread.currentThread().isInterrupted());
                }
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        synchronized (LOCK) {
            while (!waiting) {
                LOCK.wait();
            }
        }
        worker.interrupt();
        worker.join();
        System.out.println("main");
    }
}
