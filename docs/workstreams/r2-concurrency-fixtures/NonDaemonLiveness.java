public class NonDaemonLiveness {
    static volatile boolean mainPrinted;

    static class Worker extends Thread {
        public void run() {
            while (!mainPrinted) {
                Thread.onSpinWait();
            }
            System.out.println("worker");
        }
    }

    public static void main(String[] args) {
        Worker worker = new Worker();
        worker.start();
        System.out.println("main");
        mainPrinted = true;
    }
}
