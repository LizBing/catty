public class InterruptSleep {
    static volatile boolean sleeping;

    static class Worker extends Thread {
        public void run() {
            sleeping = true;
            try {
                Thread.sleep(10000);
                System.out.println("missing");
            } catch (InterruptedException expected) {
                System.out.println(Thread.currentThread().isInterrupted());
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        while (!sleeping) {
            Thread.onSpinWait();
        }
        worker.interrupt();
        worker.join();
        System.out.println("main");
    }
}
