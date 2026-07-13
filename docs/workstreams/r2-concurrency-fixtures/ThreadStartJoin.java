public class ThreadStartJoin {
    static int value;

    static class Worker extends Thread {
        public void run() {
            value = 42;
        }
    }

    public static void main(String[] args) throws Exception {
        Worker worker = new Worker();
        System.out.println(worker.isAlive());
        worker.start();
        worker.join();
        System.out.println(worker.isAlive());
        System.out.println(value);
    }
}
