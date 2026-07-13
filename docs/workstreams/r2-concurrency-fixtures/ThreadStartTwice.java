public class ThreadStartTwice {
    static class Worker extends Thread {
        public void run() {}
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        worker.start();
        worker.join();
        try {
            worker.start();
            System.out.println("missing");
        } catch (IllegalThreadStateException expected) {
            System.out.println("IllegalThreadStateException");
        }
    }
}
