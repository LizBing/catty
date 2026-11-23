package app;

import api.Base;
import api.Step;
import catty.runtime.Plugins;

public class Main {
    static int call(Step step) { return step.value(); }

    public static void main(String[] args) throws Exception {
        ClassLoader loader = Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
        Step child = (Step) Class.forName("ext.Child", false, loader)
                .getMethod("create").invoke(null);
        System.out.println(call(child));
        System.out.println(child.value());
        System.out.println(((Integer) Step.class.getMethod("value").invoke(child)).intValue());
        System.out.println(new Base().value());
    }
}
