package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
        Object plugin = Class.forName("ext.Die", true, loader)
                .getMethod("create").invoke(null);
        ((DieStep) plugin).die();
        System.out.println(0);
    }
}
