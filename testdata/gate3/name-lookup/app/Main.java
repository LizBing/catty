package app;

import catty.runtime.Plugins;
import java.lang.reflect.Method;

public class Main {
    public static void main(String[] args) throws Exception {
        Probe.reset();
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), app);

        try {
            Class.forName("ext.Worker", false, app);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        try {
            Class.forName("ext.Worker", false, null);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(2);
        }

        Class<?> worker = Class.forName("ext.Worker", false, loader);
        Method magic = worker.getMethod("magic", new Class[0]);
        System.out.println(Probe.get());
        System.out.println(((Integer) magic.invoke(null)).intValue());
        System.out.println(Probe.get());

        Class<?> base = Class.forName("app.Base", false, app);
        Method image = base.getMethod("imageMagic", new Class[0]);
        System.out.println(((Integer) image.invoke(null)).intValue());
    }
}
