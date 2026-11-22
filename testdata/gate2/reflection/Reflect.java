import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

// Reflection fixture: public method lookup, dispatch, covariant returns,
// static initialization timing, argument conversion, result boxing, receiver
// validation, lookup failures, target exception identity and caller-sensitive
// access. Every marker is an int printed through System.out.println(I)V.
public final class Reflect {
    public static int leaf(int value) {
        return value + 1;
    }

    public static int takeInt(int value) {
        return value + 100;
    }

    public static long takeLong(long value) {
        return value + 200L;
    }

    public static int takeDouble(double value) {
        return (int) (value + 300.0);
    }

    public static int takeFloat(float value) {
        return (int) (value + 400.0f);
    }

    public static int takeChar(char value) {
        return (int) value + 500;
    }

    public static boolean takeBoolean(boolean value) {
        return value;
    }

    public static int boxInt() {
        return 31;
    }

    public static long boxLong() {
        return 32L;
    }

    public static float boxFloat() {
        return 33.5f;
    }

    public static double boxDouble() {
        return 34.5;
    }

    public static boolean boxBoolean() {
        return true;
    }

    public static char boxChar() {
        return 'C';
    }

    public static short boxShort() {
        return (short) 35;
    }

    public static byte boxByte() {
        return (byte) 36;
    }

    public static void boxVoid() {
    }

    public static void main(String[] args) throws Exception {
        // Inherited lookup: Child's mirror finds Base's public method.
        Child child = new Child();
        Method inherited = Child.class.getMethod("inherited", (Class<?>[]) null);
        System.out.println(1);
        System.out.println(((Integer) inherited.invoke(child, (Object[]) null)).intValue());

        // Override lookup: Base's declaration dispatches on the actual receiver.
        Method overridden = Base.class.getMethod("overridden", (Class<?>[]) null);
        System.out.println(2);
        System.out.println(((Integer) overridden.invoke(child, (Object[]) null)).intValue());
        System.out.println(((Integer) overridden.invoke(new Base(), (Object[]) null)).intValue());

        // Covariant return: Child's mirror selects the Integer-returning override.
        Method childValue = Child.class.getMethod("value", (Class<?>[]) null);
        Method baseValue = Base.class.getMethod("value", (Class<?>[]) null);
        System.out.println(3);
        Object fromChild = childValue.invoke(child, (Object[]) null);
        System.out.println(fromChild instanceof Integer ? 1 : 0);
        System.out.println(((Integer) fromChild).intValue());
        Object fromBase = baseValue.invoke(new Base(), (Object[]) null);
        System.out.println(fromBase instanceof Integer ? 1 : 0);
        System.out.println(((Integer) fromBase).intValue());
        try {
            childValue.invoke(new Base(), (Object[]) null);
            System.out.println(-4);
        } catch (IllegalArgumentException expected) {
            System.out.println(4);
        }

        // Static method: getMethod must not initialize Init; invoke must.
        Method stats = Init.class.getMethod("stats", (Class<?>[]) null);
        System.out.println(5);
        System.out.println(((Integer) stats.invoke(null, (Object[]) null)).intValue());
        System.out.println(((Integer) stats.invoke(new Init(), (Object[]) null)).intValue());

        // Primitive widening through Method.invoke.
        Method mInt = Reflect.class.getMethod("takeInt", new Class<?>[]{Integer.TYPE});
        Method mLong = Reflect.class.getMethod("takeLong", new Class<?>[]{Long.TYPE});
        Method mDouble = Reflect.class.getMethod("takeDouble", new Class<?>[]{Double.TYPE});
        Method mFloat = Reflect.class.getMethod("takeFloat", new Class<?>[]{Float.TYPE});
        Method mChar = Reflect.class.getMethod("takeChar", new Class<?>[]{Character.TYPE});
        Method mBoolean = Reflect.class.getMethod("takeBoolean", new Class<?>[]{Boolean.TYPE});
        System.out.println(6);
        System.out.println(Integer.TYPE == int.class ? 1 : 0);
        System.out.println(((Integer) mInt.invoke(null, new Object[]{Byte.valueOf((byte) 5)})).intValue());
        System.out.println(((Integer) mInt.invoke(null, new Object[]{Short.valueOf((short) 6)})).intValue());
        System.out.println(((Integer) mInt.invoke(null, new Object[]{Character.valueOf('A')})).intValue());
        System.out.println((int) ((Long) mLong.invoke(null, new Object[]{Integer.valueOf(7)})).longValue());
        System.out.println(((Integer) mDouble.invoke(null, new Object[]{Long.valueOf(8L)})).intValue());
        System.out.println(((Integer) mDouble.invoke(null, new Object[]{Float.valueOf(1.5f)})).intValue());
        System.out.println(((Integer) mFloat.invoke(null, new Object[]{Integer.valueOf(9)})).intValue());
        System.out.println(((Integer) mDouble.invoke(null, new Object[]{Double.valueOf(10.0)})).intValue());
        System.out.println(((Integer) mFloat.invoke(null, new Object[]{Float.valueOf(1.0f)})).intValue());
        System.out.println(((Integer) mChar.invoke(null, new Object[]{Character.valueOf('B')})).intValue());
        System.out.println(((Boolean) mBoolean.invoke(null, new Object[]{Boolean.valueOf(true)})).booleanValue() ? 1 : 0);

        // Rejected conversions.
        System.out.println(7);
        try {
            mInt.invoke(null, new Object[]{Long.valueOf(1L)});
            System.out.println(-8);
        } catch (IllegalArgumentException expected) {
            System.out.println(8);
        }
        try {
            mLong.invoke(null, new Object[]{Double.valueOf(1.0)});
            System.out.println(-9);
        } catch (IllegalArgumentException expected) {
            System.out.println(9);
        }
        try {
            mInt.invoke(null, new Object[]{null});
            System.out.println(-10);
        } catch (IllegalArgumentException expected) {
            System.out.println(10);
        }
        try {
            mBoolean.invoke(null, new Object[]{Integer.valueOf(1)});
            System.out.println(-11);
        } catch (IllegalArgumentException expected) {
            System.out.println(11);
        }
        try {
            mChar.invoke(null, new Object[]{Integer.valueOf(65)});
            System.out.println(-12);
        } catch (IllegalArgumentException expected) {
            System.out.println(12);
        }

        // Primitive results come back boxed; void returns null.
        Method boxInt = Reflect.class.getMethod("boxInt", (Class<?>[]) null);
        Method boxLong = Reflect.class.getMethod("boxLong", (Class<?>[]) null);
        Method boxFloat = Reflect.class.getMethod("boxFloat", (Class<?>[]) null);
        Method boxDouble = Reflect.class.getMethod("boxDouble", (Class<?>[]) null);
        Method boxBoolean = Reflect.class.getMethod("boxBoolean", (Class<?>[]) null);
        Method boxChar = Reflect.class.getMethod("boxChar", (Class<?>[]) null);
        Method boxShort = Reflect.class.getMethod("boxShort", (Class<?>[]) null);
        Method boxByte = Reflect.class.getMethod("boxByte", (Class<?>[]) null);
        Method boxVoid = Reflect.class.getMethod("boxVoid", (Class<?>[]) null);
        System.out.println(13);
        Object boxed = boxInt.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Integer ? 1 : 0);
        System.out.println(((Integer) boxed).intValue());
        boxed = boxLong.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Long ? 1 : 0);
        System.out.println((int) ((Long) boxed).longValue());
        boxed = boxFloat.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Float ? 1 : 0);
        System.out.println((int) ((Float) boxed).floatValue());
        boxed = boxDouble.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Double ? 1 : 0);
        System.out.println((int) ((Double) boxed).doubleValue());
        boxed = boxBoolean.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Boolean ? 1 : 0);
        System.out.println(((Boolean) boxed).booleanValue() ? 1 : 0);
        boxed = boxChar.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Character ? 1 : 0);
        System.out.println((int) ((Character) boxed).charValue());
        boxed = boxShort.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Short ? 1 : 0);
        System.out.println((int) ((Short) boxed).shortValue());
        boxed = boxByte.invoke(null, (Object[]) null);
        System.out.println(boxed instanceof Byte ? 1 : 0);
        System.out.println((int) ((Byte) boxed).byteValue());
        boxed = boxVoid.invoke(null, (Object[]) null);
        System.out.println(boxed == null ? 1 : 0);

        // Receiver validation and interface dispatch.
        System.out.println(14);
        try {
            inherited.invoke(null, (Object[]) null);
            System.out.println(-15);
        } catch (NullPointerException expected) {
            System.out.println(15);
        }
        try {
            inherited.invoke(new Object(), (Object[]) null);
            System.out.println(-16);
        } catch (IllegalArgumentException expected) {
            System.out.println(16);
        }
        Method faceCall = Face.class.getMethod("call", (Class<?>[]) null);
        System.out.println(((Integer) faceCall.invoke(new FaceImpl(), (Object[]) null)).intValue());

        // Lookup failures.
        System.out.println(17);
        try {
            Child.class.getMethod("value", new Class<?>[]{Integer.TYPE});
            System.out.println(-18);
        } catch (NoSuchMethodException expected) {
            System.out.println(18);
        }
        try {
            Reflect.class.getMethod("missing", (Class<?>[]) null);
            System.out.println(-19);
        } catch (NoSuchMethodException expected) {
            System.out.println(19);
        }
        try {
            Base.class.getMethod("packagePrivate", (Class<?>[]) null);
            System.out.println(-20);
        } catch (NoSuchMethodException expected) {
            System.out.println(20);
        }

        // Target exception identity: the wrapper keeps the saved cause.
        Method explode = Thrower.class.getMethod("explode", (Class<?>[]) null);
        System.out.println(21);
        try {
            explode.invoke(new Thrower(), (Object[]) null);
            System.out.println(-22);
        } catch (InvocationTargetException wrapped) {
            System.out.println((Object) wrapped != Thrower.SAVED ? 1 : 0);
            System.out.println(wrapped.getCause() == Thrower.SAVED ? 1 : 0);
            System.out.println(wrapped.getCause() != wrapped ? 1 : 0);
        }

        // Caller-sensitive access: the declaring class is not public.
        Object secret = hidden.Factory.make();
        Method secretValue = secret.getClass().getMethod("value", (Class<?>[]) null);
        System.out.println(22);
        try {
            secretValue.invoke(secret, (Object[]) null);
            System.out.println(-23);
        } catch (IllegalAccessException expected) {
            System.out.println(23);
        }

        // A small leaf method invoked twice (forced-interpreter candidate).
        Method leafMethod = Reflect.class.getMethod("leaf", new Class<?>[]{Integer.TYPE});
        System.out.println(24);
        System.out.println(((Integer) leafMethod.invoke(null, new Object[]{Integer.valueOf(1)})).intValue());
        System.out.println(((Integer) leafMethod.invoke(null, new Object[]{Integer.valueOf(2)})).intValue());
    }
}

class Base {
    public int inherited() {
        return 11;
    }

    public int overridden() {
        return 12;
    }

    public Object value() {
        return Integer.valueOf(-1);
    }

    public int unused() {
        return 98;
    }

    int packagePrivate() {
        return 99;
    }
}

class Child extends Base {
    @Override
    public int overridden() {
        return 22;
    }

    @Override
    public Integer value() {
        return Integer.valueOf(123);
    }
}

class Init {
    static int state;

    static {
        System.out.println(30);
        state = 14;
    }

    public static int stats() {
        return state;
    }
}

class Thrower {
    static final RuntimeException SAVED = new RuntimeException();

    public void explode() {
        throw SAVED;
    }
}

interface Face {
    int call();
}

class FaceImpl implements Face {
    public int call() {
        return 41;
    }
}
