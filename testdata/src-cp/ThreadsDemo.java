public class ThreadsDemo {
    static int count = 0;
    static final Object lock = new Object();

    static void bump(int n) {
        for (int i = 0; i < n; i++) {
            synchronized (lock) { count++; }
        }
    }

    public static void main(String[] args) throws InterruptedException {
        Thread w1 = new Worker("w1", 300);
        Thread w2 = new Worker("w2", 200);
        w1.start(); w2.start();
        w1.join(); w2.join();
        System.out.println(count);

        Thread sleeper = new Sleeper();
        sleeper.start();
        Thread.sleep(50);
        sleeper.interrupt();
        sleeper.join();
        System.out.println("main done " + Thread.currentThread().getName());
        System.out.println(w1.getName() + " alive=" + w1.isAlive());
    }
}

class Worker extends Thread {
    int n;
    Worker(String name, int n) { super(name); this.n = n; }
    public void run() { ThreadsDemo.bump(n); }
}

class Sleeper extends Thread {
    public void run() {
        try {
            Thread.sleep(60_000);
            System.out.println("BUG: not interrupted");
        } catch (InterruptedException e) {
            System.out.println("interrupted ok");
        }
    }
}
