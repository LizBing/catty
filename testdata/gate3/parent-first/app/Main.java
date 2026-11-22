package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        Box box = new Box(0);
        ClassLoader app = Main.class.getClassLoader();

        ClassLoader shadow = Plugins.register(Plugins.readJar(args[0]), app);
        System.out.println(Class.forName("app.Base", false, shadow) == Base.class);
        System.out.println(Class.forName("app.Step", false, shadow) == Step.class);

        Class<?> unique = Class.forName("ext.Unique", false, shadow);
        Step impl = (Step) unique.getMethod("create").invoke(null);
        System.out.println(impl.value(box));

        ClassLoader grand = Plugins.register(Plugins.readJar(args[1]), shadow);
        Class<?> child = Class.forName("ext.Child", false, grand);
        Base childImpl = (Base) child.getMethod("create").invoke(null);
        System.out.println(childImpl.value(box));

        ClassLoader independent = Plugins.register(Plugins.readJar(args[2]), null);
        System.out.println(Class.forName("java.lang.Object", false, independent) == Object.class);
        System.out.println(Class.forName("ext.Independent", false, independent) != null);
        try {
            Class.forName("app.Base", false, independent);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        try {
            Class.forName("ext.Child", false, independent);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        try {
            Class.forName("ext.Independent", false, app);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        try {
            Class.forName("ext.Independent", false, shadow);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        System.out.println(box.get());
    }
}
