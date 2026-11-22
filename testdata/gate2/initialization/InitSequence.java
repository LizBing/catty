public final class InitSequence {
    static final Signal signal = new Signal();
    static int failures;

    static int mark(int step) {
        System.out.println(step);
        return step;
    }

    static class Parent {
        static int value = mark(10);
    }

    interface WithDefault {
        int MARKER = mark(20);

        default int marker() {
            return MARKER;
        }
    }

    interface WithoutDefault {
        int MARKER = mark(90);
    }

    static class Child extends Parent implements WithDefault, WithoutDefault {
        static int value = mark(30);

        static int read() {
            return value;
        }
    }

    static final class Signal extends RuntimeException {}

    static class Broken {
        static int value = fail();

        static int fail() {
            failures++;
            throw signal;
        }

        static int read() {
            return value;
        }
    }

    public static void main(String[] args) {
        System.out.println(Child.class == Child.class ? 1 : -1);
        System.out.println(Child.read());
        try {
            System.out.println(Broken.read());
        } catch (ExceptionInInitializerError first) {
            System.out.println(first.getCause() == signal ? 40 : -40);
        }
        try {
            System.out.println(Broken.read());
        } catch (NoClassDefFoundError later) {
            System.out.println(50);
        }
        System.out.println(failures);
    }
}
