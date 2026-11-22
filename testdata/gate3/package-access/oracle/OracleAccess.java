package app;

import catty.runtime.Plugins;

public class OracleAccess {
    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), OracleAccess.class.getClassLoader());
        try {
            Class.forName(args[1], true, loader);
            System.out.println("loaded");
        } catch (IllegalAccessError e) {
            System.out.println("IllegalAccessError");
        }
    }
}
