package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        Box box = new Box(0);
        ClassLoader good = Plugins.register(Plugins.readJar(args[0]), app);
        Base base = (Base) Class.forName("ext.Worker", false, good).getMethod("create").invoke(null);
        System.out.println(base.value(box));

        try {
            Plugins.register(Plugins.readJar(args[1]), app);
            System.out.println("unexpected");
        } catch (PluginAdmissionException e) {
            System.out.println("CLASS_FORMAT".equals(e.getCode()));
        }
        System.out.println(base.value(box));
        try {
            Class.forName("ext.Alpha", false, app);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }
        try {
            Class.forName("ext.Alpha", false, good);
            System.out.println(0);
        } catch (ClassNotFoundException e) {
            System.out.println(1);
        }

        ClassLoader corrected = Plugins.register(Plugins.readJar(args[2]), app);
        Object alpha = Class.forName("ext.Alpha", false, corrected).getMethod("create").invoke(null);
        Object result = alpha.getClass().getMethod("value", new Class[]{Box.class}).invoke(alpha, box);
        System.out.println(((Integer) result).intValue());
    }
}
