public final class InitOrder {
    static int mark(int step) {
        System.out.println(step);
        return step;
    }

    static int fail(int marker) {
        throw new InitError(marker);
    }

    static void reflect() {
        int ignored = Reflect.value;
        System.out.println(401);
    }

    static final class InitError extends Error {
        private static final long serialVersionUID = 1L;
        final int marker;

        InitError(int marker) {
            this.marker = marker;
        }
    }

    static final class RuntimeFailure extends RuntimeException {
        private static final long serialVersionUID = 1L;
    }

    static final class Signals {
        static final RuntimeFailure SAVED = new RuntimeFailure();
    }

    static final class Literal {
        static int value = mark(90);

        static int read() {
            return value;
        }
    }

    static final class Get {
        static int value = mark(100);
    }

    static final class Put {
        static int value = mark(200);
    }

    static final class Invoke {
        static int value = mark(300);

        static int touch() {
            return value + 1;
        }
    }

    static final class Reflect {
        static int value = mark(400);
    }

    static class Base {
        static int base = mark(10);

        static int readBase() {
            return base;
        }
    }

    static final class Derived extends Base {
        static int derived = mark(20);

        static int readDerived() {
            return derived;
        }
    }

    interface FaceOne {
        int BASE = mark(30);

        default int faceOne() {
            return BASE;
        }
    }

    interface FaceTwo {
        int EXTRA = mark(40);

        default int faceTwo() {
            return EXTRA;
        }
    }

    interface NoDefault {
        int MARKER = mark(900);
    }

    static final class WithFaces implements FaceOne, FaceTwo, NoDefault {
        static int value = mark(50);
        static int fromFace = FaceOne.BASE + 100;

        static int read() {
            return value;
        }
    }

    static final class ReentryA {
        static int first = mark(60);
        static int second;

        static {
            new ReentryB();
        }

        static int third = mark(64);

        static int read() {
            return third;
        }
    }

    static final class ReentryB {
        static {
            mark(61);
            mark(600 + ReentryA.first);
            mark(700 + ReentryA.second);
        }
    }

    static final class SelfA {
        static int runs;
        static int value = mark(70);
        static int reentered;

        static {
            runs = 1;
            reentered = SelfB.touch();
        }

        static int probe() {
            return value;
        }
    }

    static final class SelfB {
        static int seen;

        static {
            mark(71);
            seen = SelfA.value;
        }

        static int touch() {
            return seen;
        }
    }

    static final class Defaults {
        static int value;

        static {
            mark(value);
            value = mark(45);
        }
    }

    static final class Constants {
        static final int SCALE = 7;
        static final int MASK = 0x0F0F0F0F;
        static final long BIG = 0x1122334455667788L;
    }

    static final class BrokenError {
        static int value = fail();

        static int fail() {
            return InitOrder.fail(501);
        }

        static int read() {
            return value;
        }
    }

    static final class BrokenRuntime {
        static int value = fail();

        static int fail() {
            throw Signals.SAVED;
        }

        static int read() {
            return value;
        }
    }

    static class ParentFail {
        static int value = fail();

        static int fail() {
            return InitOrder.fail(801);
        }

        static int read() {
            return value;
        }
    }

    static final class ChildFail extends ParentFail {
        static int value = mark(850);

        static int read() {
            return value;
        }
    }

    public static void main(String[] args) {
        Class<?> literalClass = Literal.class;
        Literal[] literalArray = new Literal[0];
        if (literalClass != null && literalArray.length == 0) {
            System.out.println(1);
        }
        System.out.println(Literal.read() + 1);

        System.out.println(Get.value + 1);
        Put.value = 201;
        System.out.println(Put.value);
        System.out.println(Invoke.touch());
        reflect();

        System.out.println(Derived.readDerived() + Base.readBase());
        System.out.println(WithFaces.read() + WithFaces.fromFace);

        System.out.println(ReentryA.read() + 1);
        System.out.println(SelfA.probe());
        System.out.println(SelfA.runs);
        System.out.println(SelfA.reentered);
        System.out.println(Defaults.value + 1);

        long big = Constants.BIG;
        int high = (int) (big >>> 32);
        System.out.println(high / Constants.SCALE + Constants.MASK);

        try {
            System.out.println(BrokenError.read());
        } catch (InitError error) {
            System.out.println(error.marker);
        } catch (ExceptionInInitializerError wrapped) {
            System.out.println(-1);
        }
        try {
            BrokenError.read();
        } catch (NoClassDefFoundError later) {
            System.out.println(502);
        }

        try {
            System.out.println(BrokenRuntime.read());
        } catch (ExceptionInInitializerError wrapped) {
            System.out.println(wrapped.getCause() == Signals.SAVED ? 1 : 0);
        }
        try {
            BrokenRuntime.read();
        } catch (NoClassDefFoundError later) {
            System.out.println(602);
        }

        try {
            System.out.println(ChildFail.read());
        } catch (InitError error) {
            System.out.println(error.marker);
        } catch (ExceptionInInitializerError wrapped) {
            System.out.println(-2);
        }
        try {
            ChildFail.read();
        } catch (NoClassDefFoundError later) {
            System.out.println(802);
        }
    }
}
