package app;

import catty.runtime.Plugins;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        byte[] bytes = Plugins.readJar(args[0]);
        ClassLoader l1 = Plugins.register(bytes, app);
        ClassLoader l2 = Plugins.register(bytes, app);
        System.out.println(l1 == l2);
        Class<?> c1 = Class.forName("ext.Worker", false, l1);
        Class<?> c2 = Class.forName("ext.Worker", false, l2);
        System.out.println(c1 == c2);
        System.out.println(Class.forName("app.Base", false, l1) == Base.class);
        System.out.println(Class.forName("[I", false, l1) == Class.forName("[I", false, null));
        System.out.println(Class.forName("[Lext.Worker;", false, l1) == Class.forName("[Lext.Worker;", false, l2));

        Base w1 = (Base) c1.getMethod("create").invoke(null);
        Base w2 = (Base) c2.getMethod("create").invoke(null);
        System.out.println(w1.value(new Box(0)));
        System.out.println(w2.value(new Box(0)));
        System.out.println(w1.value(new Box(0)));
        System.out.println(w1 instanceof Step);

        Method self = c1.getMethod("self", new Class[]{Object.class});
        System.out.println(((Boolean) self.invoke(null, w1)).booleanValue() ? 1 : 0);
        System.out.println(((Boolean) self.invoke(null, w2)).booleanValue() ? 1 : 0);

        Method cast = c1.getMethod("cast", new Class[]{Object.class});
        try {
            cast.invoke(null, w2);
            System.out.println(0);
        } catch (InvocationTargetException e) {
            System.out.println(e.getCause() instanceof ClassCastException ? 3 : 9);
        }
        Method store = c1.getMethod("store", new Class[]{Object.class});
        try {
            store.invoke(null, w2);
            System.out.println(0);
        } catch (InvocationTargetException e) {
            System.out.println(e.getCause() instanceof ArrayStoreException ? 4 : 9);
        }
    }
}
