public final class State {
    static class Parent {
        int value;
    }

    static final class Child extends Parent {
        int value;
    }

    static final class Counter {
        static int total = initialize();

        static int initialize() {
            System.out.println(70);
            return 10;
        }

        static int increment() {
            total = total + 1;
            return total;
        }
    }

    static final class Fault extends RuntimeException {}

    static final class Broken {
        static int value = fail();

        static int fail() {
            throw new Fault();
        }
    }

    public static void main(String[] args) {
        Class<?> mirror = Counter.class;
        if (mirror == Counter.class) System.out.println(60);
        System.out.println(Counter.increment());
        Counter.total = Counter.total + 5;
        System.out.println(Counter.increment());

        Child child = new Child();
        Parent parent = child;
        parent.value = 21;
        child.value = 32;
        System.out.println(parent.value);
        System.out.println(child.value);

        try {
            System.out.println(Broken.value);
        } catch (ExceptionInInitializerError expected) {
            System.out.println(81);
        }
        try {
            System.out.println(Broken.value);
        } catch (NoClassDefFoundError expected) {
            System.out.println(82);
        }
    }
}
