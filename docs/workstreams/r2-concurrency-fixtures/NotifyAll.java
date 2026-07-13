public class NotifyAll {
    static final Object LOCK = new Object();
    static int waiting;
    static int resumed;
    static boolean released;

    static class Worker extends Thread {
        public void run() {
            synchronized (LOCK) {
                waiting++;
                LOCK.notifyAll();
                while (!released) {
                    try {
                        LOCK.wait();
                    } catch (InterruptedException unexpected) {
                        throw new RuntimeException(unexpected);
                    }
                }
                resumed++;
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Worker first = new Worker();
        Worker second = new Worker();
        first.start();
        second.start();
        synchronized (LOCK) {
            while (waiting != 2) {
                LOCK.wait();
            }
            released = true;
            LOCK.notifyAll();
        }
        first.join();
        second.join();
        System.out.println(resumed);
    }
}
