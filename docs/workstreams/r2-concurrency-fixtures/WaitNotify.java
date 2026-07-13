public class WaitNotify {
    static final Object LOCK = new Object();
    static boolean waiting;
    static boolean released;

    static class Worker extends Thread {
        public void run() {
            synchronized (LOCK) {
                waiting = true;
                LOCK.notifyAll();
                while (!released) {
                    try {
                        LOCK.wait();
                    } catch (InterruptedException unexpected) {
                        throw new RuntimeException(unexpected);
                    }
                }
            }
            System.out.println("worker");
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        synchronized (LOCK) {
            while (!waiting) {
                LOCK.wait();
            }
            released = true;
            LOCK.notifyAll();
        }
        worker.join();
        System.out.println("main");
    }
}
