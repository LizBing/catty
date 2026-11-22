package app;

import catty.runtime.Plugins;
import demo.Base;
import demo.PublicBase;

public class Main {
    static int baseCall(Base base) { return base.call(); }

    static int publicCall(PublicBase base) { return base.call(); }

    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
        Base imageBase = new Base();
        Base child = (Base) Class.forName("demo.Child", false, loader).getMethod("create").invoke(null);
        PublicBase imagePublic = new PublicBase();
        PublicBase childPublic = (PublicBase) Class.forName("demo.PublicChild", false, loader)
                .getMethod("create").invoke(null);

        System.out.println(baseCall(imageBase));
        System.out.println(baseCall(child));
        System.out.println(publicCall(imagePublic));
        System.out.println(publicCall(childPublic));
    }
}
