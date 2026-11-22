package app;

import catty.runtime.Plugins;

public class OracleMain {
    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), OracleMain.class.getClassLoader());
        try {
            Class.forName("app.ProbeBad", true, loader);
            System.out.println("loaded");
        } catch (VerifyError e) {
            System.out.println("VerifyError");
        }
    }
}
