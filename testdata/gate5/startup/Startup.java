/** Startup arguments, class initialization order and one mixed engine crossing. */
public final class Startup {
    static {
        System.out.println("clinit");
    }

    static final class Box {
        int count;

        void bump() {
            count++;
        }
    }

    // Forced to interpretation in the mixed image; the caller and callback stay AOT.
    static void step(Box box) {
        box.bump();
        box.bump();
    }

    public static void main(String[] args) {
        System.out.println("main");
        System.out.println(args.length);
        for (int i = 0; i < args.length; i++) {
            System.out.println(args[i]);
        }
        Box box = new Box();
        step(box);
        System.out.println(box.count);
    }
}
