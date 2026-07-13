public class CrossThreadClassInitialization {
    static class Target {
        static int initializationCount;
        static final int VALUE;

        static {
            initializationCount++;
            try {
                Thread.sleep(50);
            } catch (InterruptedException unexpected) {
                throw new RuntimeException(unexpected);
            }
            VALUE = 7;
        }
    }

    static int firstResult;
    static int secondResult;

    static class First extends Thread {
        public void run() {
            firstResult = Target.VALUE;
        }
    }

    static class Second extends Thread {
        public void run() {
            secondResult = Target.VALUE;
        }
    }

    public static void main(String[] args) throws Exception {
        First first = new First();
        Second second = new Second();
        first.start();
        second.start();
        first.join();
        second.join();
        System.out.println(Target.initializationCount);
        System.out.println(firstResult);
        System.out.println(secondResult);
    }
}
