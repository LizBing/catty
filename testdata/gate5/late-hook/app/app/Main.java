package app;

import catty.runtime.Plugins;
import java.lang.reflect.Method;

/** Built before the plugin JAR exists; registers a plugin-backed shutdown hook. */
public final class Main {
    public static void main(String[] args) throws Exception {
        byte[] jar = Plugins.readJar(args[0]);
        ClassLoader parent = Main.class.getClassLoader();
        ClassLoader loader = Plugins.register(jar, parent);
        Class<?> factory = Class.forName("ext.HookFactory", true, loader);
        Method make = factory.getMethod("make", new Class[0]);
        Thread hook = (Thread) make.invoke(null, new Object[0]);
        Runtime.getRuntime().addShutdownHook(hook);
        System.out.println("registered");
        System.exit(5);
    }
}
