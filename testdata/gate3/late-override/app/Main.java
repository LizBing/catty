package app;

import catty.runtime.Plugins;

public class Main {
    static int virt(Base base, Box box) { return base.value(box); }

    static int inter(Step step, Box box) { return step.value(box); }

    public static void main(String[] args) throws Exception {
        Box box = new Box(0);
        Base base = new Base();
        Step step = base;

        System.out.println(virt(base, box));
        System.out.println(virt(base, box));
        System.out.println(inter(step, box));
        System.out.println(inter(step, box));
        System.out.println(1);

        ClassLoader a = Plugins.register(Plugins.readJar(args[0]), Main.class.getClassLoader());
        Base lateA = (Base) Class.forName("ext.Worker", false, a).getMethod("create").invoke(null);
        System.out.println(virt(lateA, box));
        System.out.println(virt(lateA, box));
        System.out.println(inter(lateA, box));
        System.out.println(inter(lateA, box));

        ClassLoader b = Plugins.register(Plugins.readJar(args[1]), Main.class.getClassLoader());
        Base lateB = (Base) Class.forName("ext.Worker", false, b).getMethod("create").invoke(null);
        System.out.println(virt(lateB, box));
        System.out.println(virt(lateA, box));
        System.out.println(box.get());
        System.out.println(box.callbacks());
    }
}
