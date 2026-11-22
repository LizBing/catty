package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    static int call(ClassLoader loader) throws Exception {
        Object value = Class.forName("good.Handler", false, loader)
                .getMethod("f", int.class).invoke(null, 0);
        return ((Integer) value).intValue();
    }

    static boolean rejected(String path, ClassLoader app) throws Exception {
        try {
            Plugins.register(Plugins.readJar(path), app);
            return false;
        } catch (PluginAdmissionException e) {
            return "VERIFY".equals(e.getCode());
        }
    }

    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader good = Plugins.register(Plugins.readJar(args[0]), app);
        System.out.println(call(good));
        System.out.println(rejected(args[1], app));
        System.out.println(rejected(args[2], app));
        ClassLoader again = Plugins.register(Plugins.readJar(args[0]), app);
        System.out.println(call(again));
    }
}
