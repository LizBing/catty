public class VolatilePublication {
    static int value;
    static volatile boolean ready;
    static int observed;

    static class Worker extends Thread {
        public void run() {
            while (!ready) {
                Thread.onSpinWait();
            }
            observed = value;
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        value = 42;
        ready = true;
        worker.join();
        System.out.println(observed);
    }
}
