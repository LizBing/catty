import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public class Monitors {
    private static final RuntimeException TOKEN = new RuntimeException();
    private static final Object BLOCK = new Object();
    private static final Monitors BOX = new Monitors();
    private static volatile boolean entering;
    private static boolean blockedInterrupt;
    private static int count;

    public synchronized void fail() {
        Check.that(Thread.holdsLock(this), "instance method monitor");
        nested();
    }

    private synchronized void nested() {
        synchronized (this) {
            Check.that(Thread.holdsLock(this), "reentrant block monitor");
            throw TOKEN;
        }
    }

    public static synchronized void failStatic() {
        Check.that(Thread.holdsLock(Monitors.class), "static method monitor");
        throw TOKEN;
    }

    private static void block(Object monitor) {
        synchronized (monitor) {
            count++;
        }
    }

    static class Worker extends Thread {
        public void run() {
            entering = true;
            synchronized (BLOCK) {
                blockedInterrupt = Thread.interrupted();
                for (int i = 0; i < 500; i++) {
                    block(BLOCK);
                }
            }
            synchronized (BOX) {
                synchronized (Monitors.class) {
                    Check.that(Thread.holdsLock(BOX), "worker could not acquire instance monitor");
                }
            }
        }
    }

    public static void main(String[] args) throws Exception {
        synchronized (BOX) {
            try {
                BOX.fail();
                throw new RuntimeException("missing instance throw");
            } catch (RuntimeException caught) {
                Check.that(caught == TOKEN && Thread.holdsLock(BOX), "unwind dropped the enclosing lock or replaced throwable");
            }
        }
        Method method = Monitors.class.getMethod("failStatic", new Class[0]);
        try {
            method.invoke(null, new Object[0]);
            throw new RuntimeException("missing reflective throw");
        } catch (InvocationTargetException caught) {
            Check.that(caught.getCause() == TOKEN, "reflection replaced throw cause");
        }
        Check.that(!Thread.holdsLock(BOX) && !Thread.holdsLock(Monitors.class), "method monitor leaked");
        Worker worker = new Worker();
        synchronized (BLOCK) {
            worker.start();
            while (!entering) {
                Thread.yield();
            }
            worker.interrupt();
            Check.that(count == 0, "contender entered too early");
        }
        for (int i = 0; i < 500; i++) {
            block(BLOCK);
        }
        worker.join();
        Check.that(blockedInterrupt && count == 1000, "monitor entry ignored ownership or swallowed interrupt");
        try {
            block(null);
            throw new RuntimeException("null monitor entered");
        } catch (NullPointerException expected) {
            System.out.println("monitors");
        }
    }
}
