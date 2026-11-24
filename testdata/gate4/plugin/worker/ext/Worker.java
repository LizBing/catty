package ext;

public class Worker extends Thread {
    private static volatile int visits;

    private static synchronized void record() {
        if (!Thread.holdsLock(Worker.class)) {
            throw new RuntimeException("plugin static method used the wrong monitor");
        }
        visits++;
    }

    public void run() {
        if (this != Thread.currentThread()) {
            throw new RuntimeException("plugin thread has the wrong identity");
        }
        record();
    }

    public static Thread make() {
        return new Worker();
    }

    public static int count() {
        return visits;
    }
}
