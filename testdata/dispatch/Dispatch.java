import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public final class Dispatch {
    public interface Choice {
        int pick();
    }

    public static class Parent {
        public Object reveal(Token token) {
            return token;
        }

        public Object reveal(Other other) {
            return null;
        }
    }

    public static final class First extends Parent implements Choice {
        public int pick() {
            return 11;
        }

        public Object reveal(Token token) {
            return null;
        }
    }

    public static final class Second extends Parent implements Choice {
        public int pick() {
            return 22;
        }
    }

    public static final class Token {}
    public static final class Other {}
    public static final class Alarm extends RuntimeException {}

    static int choose(Choice receiver) {
        return receiver.pick();
    }

    public static int increment(int n) {
        return n + 3;
    }

    public static void fail(Alarm alarm) {
        throw alarm;
    }

    public static void main(String[] args) throws ReflectiveOperationException {
        First first = new First();
        Second second = new Second();
        Choice receiver = first;
        Choice alternate = second;
        if (args.length == 2) {
            receiver = second;
            alternate = first;
        }
        System.out.println(choose(receiver));
        System.out.println(choose(alternate));

        Token token = new Token();
        Method inherited = Second.class.getMethod(args[0], Token.class);
        if (inherited.invoke(second, token) == token) System.out.println(31);
        Object nil = null;
        if (inherited.invoke(first, token) == nil) System.out.println(32);

        try {
            Second.class.getMethod("absent", Token.class);
        } catch (NoSuchMethodException expected) {
            System.out.println(41);
        }
        try {
            inherited.invoke(second, new Other());
        } catch (IllegalArgumentException expected) {
            System.out.println(42);
        }

        Method increment = Dispatch.class.getMethod("increment", Integer.TYPE);
        Object boxed = increment.invoke(null, Integer.valueOf(9));
        if (boxed == Integer.valueOf(12)) System.out.println(Integer.valueOf(12).intValue());
        Class<?>[] noTypes = null;
        Object[] noArgs = null;
        Method intValue = Integer.class.getMethod("intValue", noTypes);
        if (intValue.invoke(boxed, noArgs) == boxed) System.out.println(43);

        Alarm alarm = new Alarm();
        try {
            Dispatch.class.getMethod("fail", Alarm.class).invoke(null, alarm);
        } catch (InvocationTargetException e) {
            if (e.getCause() == alarm) System.out.println(51);
        }
    }
}
