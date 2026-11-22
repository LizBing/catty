package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader full = Plugins.register(Plugins.readJar(args[0]), app);
        ClassLoader missing = Plugins.register(Plugins.readJar(args[1]), app);

        System.out.println(((Integer) Class.forName("ext.Link", false, full).getMethod("good").invoke(null)).intValue());
        System.out.println(((Integer) Class.forName("ext.Link", false, missing).getMethod("good").invoke(null)).intValue());
        System.out.println(Class.forName("ext.Sig", false, missing) != null);

        for (int i = 0; i < 2; i++) {
            try {
                Class.forName("ext.Link", false, missing).getMethod("field").invoke(null);
                System.out.println(0);
            } catch (java.lang.reflect.InvocationTargetException e) {
                System.out.println(e.getCause() instanceof NoSuchFieldError ? 1 : 9);
            }
        }
        System.out.println(((Integer) Class.forName("ext.Link", false, full).getMethod("field").invoke(null)).intValue());

        for (int i = 0; i < 2; i++) {
            try {
                Class.forName("ext.Link", false, missing).getMethod("method").invoke(null);
                System.out.println(0);
            } catch (java.lang.reflect.InvocationTargetException e) {
                System.out.println(e.getCause() instanceof NoSuchMethodError ? 2 : 9);
            }
        }
        System.out.println(((Integer) Class.forName("ext.Link", false, full).getMethod("method").invoke(null)).intValue());

        for (int i = 0; i < 2; i++) {
            try {
                Class.forName("ext.Link", false, missing).getMethod("klass").invoke(null);
                System.out.println(0);
            } catch (java.lang.reflect.InvocationTargetException e) {
                System.out.println(e.getCause() instanceof NoClassDefFoundError ? 3 : 9);
            }
        }
        System.out.println(((Integer) Class.forName("ext.Link", false, full).getMethod("good").invoke(null)).intValue());
    }
}
