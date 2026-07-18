import java.lang.reflect.Method;

public class ReflectiveConversions {
    static class Subject {
        public long widen(long left, double right) {
            return left + (long) right;
        }

        public int arrayLength(String... values) {
            return values.length;
        }
    }

    public static void main(String[] args) throws Exception {
        Subject subject = new Subject();
        Method widen = Subject.class.getMethod("widen", long.class, double.class);
        Method varargs = Subject.class.getMethod("arrayLength", String[].class);
        Object widened = widen.invoke(subject, Byte.valueOf((byte) 5), Float.valueOf(2.75f));
        Object count = varargs.invoke(subject, (Object) new String[] {"a", "b", "c"});
        System.out.println(widened);
        System.out.println(count);
        System.out.println(varargs.isVarArgs());
    }
}
