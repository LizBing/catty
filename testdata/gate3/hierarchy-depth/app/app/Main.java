package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    static int depth(ClassLoader loader, String name) throws Exception {
        return ((Integer) Class.forName(name, false, loader).getMethod("depth").invoke(null)).intValue();
    }

    static String rejected(String jar, ClassLoader parent) throws Exception {
        try {
            Plugins.register(Plugins.readJar(jar), parent);
            return "accepted";
        } catch (PluginAdmissionException e) {
            return e.getCode();
        }
    }

    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader big = Plugins.register(Plugins.readJar(args[0]), app);
        System.out.println(depth(big, "chain.C127"));
        ClassLoader p1 = Plugins.register(Plugins.readJar(args[2]), app);
        System.out.println(depth(p1, "span.P119"));
        ClassLoader p2 = Plugins.register(Plugins.readJar(args[3]), p1);
        System.out.println(depth(p2, "span.P127"));
        System.out.println(rejected(args[1], app));
        System.out.println(rejected(args[4], p2));
        System.out.println(depth(big, "chain.C127") == 128 && depth(p2, "span.P127") == 128);
    }
}
