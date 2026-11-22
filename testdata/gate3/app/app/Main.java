package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        Box box = new Box(0);
        Base base = new Base();
        Step step = base;

        System.out.println(base.value(box));
        System.out.println(base.value(box));
        System.out.println(step.value(box));
        System.out.println(step.value(box));
        System.out.println(2);

        byte[] jar = Plugins.readJar(args[0]);
        ClassLoader loader = Plugins.register(jar, Main.class.getClassLoader());
        Class<?> worker = Class.forName("ext.Worker", false, loader);
        Base late = (Base) worker.getMethod("create").invoke(null);

        System.out.println(base.value(box));
        base = late;
        step = late;
        System.out.println(base.value(box));
        System.out.println(base.value(box));
        System.out.println(step.value(box));
        System.out.println(step.value(box));
        System.out.println(box.get());
        System.out.println(box.callbacks());
        System.out.println(loader == Main.class.getClassLoader());
        System.out.println(worker.getClassLoader() == loader);
    }
}
