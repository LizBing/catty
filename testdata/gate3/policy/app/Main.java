package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        try {
            Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
            System.out.println("accepted");
        } catch (PluginAdmissionException e) {
            System.out.println(e.getCode());
        }
    }
}
