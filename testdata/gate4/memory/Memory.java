public class Memory {
    private static int startInput;
    private static int joinResult;
    private static int plain;
    private static final long[] WORDS = new long[1];
    private static final Box[] BOXES = new Box[1];
    private static volatile int ready;
    private static volatile long x;
    private static volatile long y;
    private static long readX;
    private static long readY;
    private static Box racy;

    static class Box {
        final int finalValue;
        int ordinaryValue;

        Box(int value) {
            finalValue = value;
            ordinaryValue = value + 1;
        }
    }

    static class Publisher extends Thread {
        public void run() {
            int value = startInput;
            joinResult = value + 1;
            plain = value + 2;
            WORDS[0] = 0x1122334455667788L;
            BOXES[0] = new Box(value);
            ready = 1;
        }
    }

    static class Left extends Thread {
        public void run() {
            x = 1;
            readY = y;
        }
    }

    static class Right extends Thread {
        public void run() {
            y = 1;
            readX = x;
        }
    }

    static class RacyPublisher extends Thread {
        public void run() {
            racy = new Box(19);
        }
    }

    private static void publication() throws InterruptedException {
        startInput = 40;
        Publisher publisher = new Publisher();
        publisher.start();
        while (ready == 0) {
            Thread.yield();
        }
        Box box = BOXES[0];
        Check.that(plain == 42 && WORDS[0] == 0x1122334455667788L
                && box != null && box.finalValue == 40 && box.ordinaryValue == 41,
                "volatile publication lost a field, array element or constructor write");
        publisher.join();
        Check.that(joinResult == 41, "start or join publication failed");
    }

    private static void litmus() throws InterruptedException {
        x = 0;
        y = 0;
        Left left = new Left();
        Right right = new Right();
        left.start();
        right.start();
        left.join();
        right.join();
        Check.that((readX == 0 || readX == 1) && (readY == 0 || readY == 1)
                && (readX != 0 || readY != 0), "forbidden volatile store-buffering outcome");

        racy = null;
        RacyPublisher publisher = new RacyPublisher();
        publisher.start();
        // Read before join, so the race remains visible to the host race detector.
        for (int i = 0; i < 64; i++) {
            Box seen = racy;
            if (seen != null) {
                int ordinary = seen.ordinaryValue;
                Check.that(seen.finalValue == 19 && (ordinary == 0 || ordinary == 20),
                        "unsafe constructor publication produced a forbidden value");
                break;
            }
            Thread.yield();
        }
        publisher.join();
        Check.that(racy != null && racy.ordinaryValue == 20, "join failed to publish ordinary constructor field");
    }

    public static void main(String[] args) throws InterruptedException {
        publication();
        for (int i = 0; i < 64; i++) {
            litmus();
        }
        System.out.println("memory");
    }
}
