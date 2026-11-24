package app;

import catty.runtime.Plugins;

/** Gate 4-compatible admission-boundary host: registers a late JAR and reports. */
public final class Main {
    public static void main(String[] args) throws Exception {
        byte[] jar = Plugins.readJar(args[0]);
        Plugins.register(jar, Main.class.getClassLoader());
        System.out.println("admitted");
    }
}
