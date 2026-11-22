package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    static void expectRejected(String path, ClassLoader app) throws Exception {
        try {
            Plugins.register(Plugins.readJar(path), app);
            System.out.println(0);
        } catch (PluginAdmissionException e) {
            System.out.println("VERIFY".equals(e.getCode()));
        }
    }

    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        expectRejected(args[0], app);
        expectRejected(args[1], app);
        ClassLoader good = Plugins.register(Plugins.readJar(args[2]), app);
        Object child = Class.forName("demo.LocalChild", false, good).getMethod("create").invoke(null);
        System.out.println(((Integer) child.getClass().getMethod("value").invoke(child)).intValue());
        Object impl = Class.forName("demo.LocalImpl", false, good).getMethod("create").invoke(null);
        System.out.println(((Integer) impl.getClass().getMethod("value").invoke(impl)).intValue());
    }
}
