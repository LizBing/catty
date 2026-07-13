public class SynchronizedBlocks {
    static final Object LOCK = new Object();
    static int counter;

    static void nested() {
        synchronized (LOCK) {
            counter++;
        }
    }

    static class Worker extends Thread {
        public void run() {
            for (int i = 0; i < 1000; i++) {
                synchronized (LOCK) {
                    nested();
                }
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Worker first = new Worker();
        Worker second = new Worker();
        first.start();
        second.start();
        first.join();
        second.join();
        System.out.println(counter);
    }
}
