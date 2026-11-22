package app;

import catty.runtime.Plugins;

public class OracleHandler {
    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), OracleHandler.class.getClassLoader());
        try {
            Class.forName(args[1], true, loader);
            System.out.println("loaded");
        } catch (VerifyError e) {
            System.out.println("VerifyError");
        }
    }
}
