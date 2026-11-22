package app;

import catty.runtime.Plugins;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
        Class<?> cls = Class.forName("ext.Thrower", false, loader);
        Object thrower = cls.getMethod("create").invoke(null);
        ThrowStep direct = (ThrowStep) thrower;
        Throwable original = new RuntimeException("original");
        Box box = new Box(0);
        try {
            direct.boom(box, original);
            System.out.println(0);
        } catch (Throwable caught) {
            System.out.println(caught == original);
            System.out.println(box.get());
        }
        Method reflected = cls.getMethod("boom", new Class[]{Box.class, Throwable.class});
        try {
            reflected.invoke(thrower, box, original);
            System.out.println(0);
        } catch (InvocationTargetException e) {
            System.out.println(e.getCause() == original);
            System.out.println(box.get());
        }
        System.out.println(box.callbacks());
    }
}
