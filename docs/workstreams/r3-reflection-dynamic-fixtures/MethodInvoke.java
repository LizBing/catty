import java.lang.reflect.Method;

public class MethodInvoke {
    static class Subject {
        public long combine(int left, long right) {
            return left * 10L + right;
        }

        public String virtualName() {
            return "subject";
        }
    }

    static class Child extends Subject {
        @Override
        public String virtualName() {
            return "child";
        }
    }

    public static void main(String[] args) throws Exception {
        Method combine = Subject.class.getMethod("combine", int.class, long.class);
        Method virtual = Subject.class.getMethod("virtualName");
        Object result = combine.invoke(new Subject(), Integer.valueOf(3), Long.valueOf(7));
        System.out.println(result.getClass() == Long.class);
        System.out.println(((Long) result).longValue());
        System.out.println(virtual.invoke(new Child()));
    }
}
