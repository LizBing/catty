public class DaemonLiveness {
    static class Worker extends Thread {
        public void run() {
            try {
                Thread.sleep(10000);
                System.out.println("missing");
            } catch (InterruptedException unexpected) {
                System.out.println("interrupted");
            }
        }
    }

    public static void main(String[] args) {
        Worker worker = new Worker();
        worker.setDaemon(true);
        worker.start();
        System.out.println("main");
    }
}
