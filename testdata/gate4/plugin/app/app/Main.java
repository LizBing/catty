package app;

import catty.runtime.Plugins;
import java.lang.reflect.Method;

public class Main {
    private static byte[] jar;
    private static ClassLoader parent;

    static class Registrar extends Thread {
        Class<?> workerClass;
        int visits;

        public void run() {
            try {
                ClassLoader loader = Plugins.register(jar, parent);
                workerClass = Class.forName("ext.Worker", true, loader);
                Method make = workerClass.getMethod("make", new Class[0]);
                Thread worker = (Thread) make.invoke(null, new Object[0]);
                worker.start();
                worker.join();
                Method count = workerClass.getMethod("count", new Class[0]);
                visits = ((Integer) count.invoke(null, new Object[0])).intValue();
            } catch (Exception error) {
                throw new RuntimeException("plugin failed");
            }
        }
    }

    public static void main(String[] args) throws Exception {
        jar = Plugins.readJar(args[0]);
        parent = Main.class.getClassLoader();
        Registrar a = new Registrar();
        Registrar b = new Registrar();
        a.start();
        b.start();
        a.join();
        b.join();
        if (a.workerClass == null || b.workerClass == null
                || a.workerClass == b.workerClass || a.visits != 1 || b.visits != 1) {
            throw new RuntimeException("loader-local thread or static identity was lost");
        }
        System.out.println("plugin");
    }
}
