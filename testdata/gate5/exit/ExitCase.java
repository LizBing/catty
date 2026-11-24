/** Explicit exits from workers and reflection, caller liveness, and status encoding. */
public final class ExitCase {
    static final Object GATE = new Object();
    static final Object MONITOR = new Object();
    static volatile boolean ready;
    static volatile boolean ack;
    static volatile boolean hookChecked;
    static volatile Thread exitCaller;
    static volatile boolean monitorHook;

    static int statusFor(String mode) {
        if (mode.equals("worker7") || mode.equals("reflect7")) {
            return 7;
        }
        if (mode.equals("workerneg") || mode.equals("reflectneg")) {
            return -3;
        }
        if (mode.equals("workerwide") || mode.equals("reflectwide")) {
            return 300;
        }
        return 0;
    }

    static final class Hook implements Runnable {
        public void run() {
            System.out.println("hook");
            Thread caller = exitCaller;
            if (caller != null) {
                System.out.println(caller.isAlive() ? "caller-alive=1" : "caller-alive=0");
                hookChecked = true;
            }
            if (monitorHook) {
                synchronized (MONITOR) {
                    System.out.println("acquired");
                }
            }
        }
    }

    static final class Halter implements Runnable {
        public void run() {
            while (!hookChecked) {
                Thread.yield();
            }
            try {
                Thread.sleep(40);
            } catch (InterruptedException ignored) {
            }
            System.out.println("halt-hook");
            Runtime.getRuntime().halt(23);
        }
    }

    static final class Exiter implements Runnable {
        private final int code;
        private final boolean monitor;

        Exiter(int code, boolean monitor) {
            this.code = code;
            this.monitor = monitor;
        }

        public void run() {
            exitCaller = Thread.currentThread();
            ready = true;
            while (!ack) {
                Thread.yield();
            }
            if (monitor) {
                synchronized (MONITOR) {
                    System.out.println("worker-exit");
                    try {
                        System.exit(0);
                    } finally {
                        System.out.println("finally");
                    }
                    System.out.println("after-exit");
                }
            } else {
                synchronized (GATE) {
                    System.out.println("worker-exit");
                    try {
                        System.exit(code);
                    } finally {
                        System.out.println("finally");
                    }
                    System.out.println("after-exit");
                }
            }
        }
    }

    static void reflectiveExit(Runtime runtime, int code) throws Exception {
        java.lang.reflect.Method exit = Runtime.class.getMethod("exit", new Class[]{Integer.TYPE});
        System.out.println("reflected");
        exit.invoke(runtime, new Object[]{Integer.valueOf(code)});
        System.out.println("after-reflect");
    }

    public static void main(String[] args) throws Exception {
        String mode = args[0];
        if (mode.equals("clinit")) {
            ClinitExit.main(new String[0]);
            return;
        }
        Runtime runtime = Runtime.getRuntime();
        if (mode.equals("reflect0") || mode.equals("reflect7") || mode.equals("reflectneg")
                || mode.equals("reflectwide")) {
            runtime.addShutdownHook(new Thread(new Hook(), "hook"));
            System.out.println("main-reflect");
            reflectiveExit(runtime, statusFor(mode));
        } else if (mode.equals("monitor")) {
            runtime.addShutdownHook(new Thread(new Hook(), "hook"));
            runtime.addShutdownHook(new Thread(new Halter(), "halter"));
            monitorHook = true;
            Thread worker = new Thread(new Exiter(0, true), "exit-worker");
            worker.start();
            while (!ready) {
                Thread.yield();
            }
            System.out.println("main-blocked");
            ack = true;
            for (;;) {
                Thread.yield();
            }
        } else {
            runtime.addShutdownHook(new Thread(new Hook(), "hook"));
            Thread worker = new Thread(new Exiter(statusFor(mode), false), "exit-worker");
            worker.start();
            while (!ready) {
                Thread.yield();
            }
            System.out.println("main-blocked");
            ack = true;
            synchronized (GATE) {
                try {
                    GATE.wait();
                } catch (InterruptedException ignored) {
                }
            }
        }
        System.out.println("main-end");
    }
}
