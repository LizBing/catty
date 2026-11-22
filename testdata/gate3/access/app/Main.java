package app;

import catty.runtime.PluginAdmissionException;
import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader valid = Plugins.register(Plugins.readJar(args[0]), app);
        Object probe = Class.forName("app.Probe2", false, valid).getMethod("create").invoke(null);
        Object direct = probe.getClass().getMethod("call").invoke(probe);
        System.out.println(((Integer) direct).intValue());
        Object nest = Class.forName("ext.NestAccess", false, valid).getMethod("value").invoke(null);
        System.out.println(((Integer) nest).intValue());
        Base asBase = (Base) probe;
        System.out.println(asBase.tag());
        try {
            Plugins.register(Plugins.readJar(args[1]), app);
            System.out.println(0);
        } catch (PluginAdmissionException e) {
            System.out.println("VERIFY".equals(e.getCode()));
        }
    }
}
