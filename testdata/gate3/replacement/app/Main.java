package app;

import catty.runtime.Plugins;

public class Main {
    public static void main(String[] args) throws Exception {
        Box box = new Box(0);
        ClassLoader app = Main.class.getClassLoader();
        ClassLoader la = Plugins.register(Plugins.readJar(args[0]), app);
        Base a = (Base) Class.forName("ext.Worker", false, la).getMethod("create").invoke(null);
        System.out.println(a.value(box));
        ClassLoader lb = Plugins.register(Plugins.readJar(args[1]), app);
        Base b = (Base) Class.forName("ext.Worker", false, lb).getMethod("create").invoke(null);
        System.out.println(b.value(box));
        System.out.println(a.value(box));
        System.out.println(a instanceof Base);
        System.out.println(b instanceof Base);
        System.out.println(box.get());
    }
}
