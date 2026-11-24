import catty.runtime.Plugins;
import java.lang.reflect.Method;

/** Instrumented workload exercising every Gate 5 counter family. */
public final class Telemetry {
    static final Object LOCK = new Object();
    static volatile boolean contenderReady;
    static volatile int contended;

    static final class Box {
        int value;

        void bump() {
            value++;
        }
    }

    // Forced to interpretation in the mixed image; caller and callback stay AOT.
    static void crossing(Box box) {
        box.bump();
    }

    public static int probe() {
        return 2;
    }

    static final class Contender implements Runnable {
        public void run() {
            contenderReady = true;
            synchronized (LOCK) {
                contended++;
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Runtime runtime = Runtime.getRuntime();
        ClassLoader parent = Telemetry.class.getClassLoader();

        byte[] jar = Plugins.readJar(args[0]);
        ClassLoader loader = Plugins.register(jar, parent);
        Class<?> taskClass = Class.forName("ext.PluginTask", true, loader);
        Method make = taskClass.getMethod("make", new Class[0]);
        Thread task = (Thread) make.invoke(null, new Object[0]);
        task.start();
        task.join();

        if (!args[1].equals("skip")) {
            try {
                byte[] invalid = Plugins.readJar(args[1]);
                Plugins.register(invalid, parent);
                System.out.println("accepted-invalid");
            } catch (Throwable expected) {
                System.out.println("rejected");
            }
        }

        Box box = new Box();
        crossing(box);
        System.out.println(box.value);

        Method probe = Telemetry.class.getMethod("probe", new Class[0]);
        System.out.println(((Integer) probe.invoke(null, new Object[0])).intValue());
        try {
            Telemetry.class.getMethod("absent", new Class[0]);
            System.out.println("found-absent");
        } catch (NoSuchMethodException expected) {
            System.out.println("lookup-miss");
        }

        try {
            throw new RuntimeException("matched");
        } catch (RuntimeException expected) {
            System.out.println("caught");
        }

        Thread contender = new Thread(new Contender(), "contender");
        synchronized (LOCK) {
            contender.start();
            while (!contenderReady) {
                Thread.yield();
            }
            // Keep the monitor held long enough for the contender to park on it.
            try {
                Thread.sleep(200);
            } catch (InterruptedException ignored) {
            }
        }
        contender.join();
        System.out.println(contended == 1 ? "contended" : "uncontended");

        runtime.gc();
        System.out.println("done");
    }
}
