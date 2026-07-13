public class InterruptJoin {
    static volatile boolean joining;

    static class Target extends Thread {
        public void run() {
            try {
                Thread.sleep(10000);
            } catch (InterruptedException expected) {}
        }
    }

    static class Joiner extends Thread {
        final Thread target;

        Joiner(Thread target) {
            this.target = target;
        }

        public void run() {
            joining = true;
            try {
                target.join();
                System.out.println("missing");
            } catch (InterruptedException expected) {
                System.out.println(Thread.currentThread().isInterrupted());
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Target target = new Target();
        Joiner joiner = new Joiner(target);
        target.start();
        joiner.start();
        while (!joining) {
            Thread.onSpinWait();
        }
        joiner.interrupt();
        joiner.join();
        target.interrupt();
        target.join();
        System.out.println("main");
    }
}
