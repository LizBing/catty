package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        byte[] bytes = Plugins.readJar(args[0]);
        ClassLoader a = Plugins.register(bytes, app);
        Probe.reset();
        Class<?> init = Class.forName("ext.Init", false, a);
        Object array = Class.forName("ext.Factory", false, a).getMethod("make").invoke(null);
        System.out.println(Probe.get());
        Class.forName("ext.Init", true, a);
        System.out.println(Probe.get());
        System.out.println(((Integer) init.getMethod("value").invoke(null)).intValue());

        ClassLoader b = Plugins.register(bytes, app);
        Class.forName("ext.Init", true, b);
        System.out.println(Probe.get());

        Probe.supplied = new RuntimeException("supplied");
        ClassLoader c = Plugins.register(bytes, app);
        try {
            Class.forName("ext.Bad", true, c);
            System.out.println(0);
        } catch (ExceptionInInitializerError e) {
            System.out.println(e.getCause() == Probe.supplied);
        }
        try {
            Class.forName("ext.Bad", true, c);
            System.out.println(0);
        } catch (NoClassDefFoundError e) {
            System.out.println(1);
        }
    }
}
