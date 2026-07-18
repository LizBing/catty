import java.lang.reflect.Constructor;
import java.lang.reflect.Field;
import java.lang.reflect.Method;

public class DeclaredMembers {
    static class Subject {
        public int alpha;
        private String beta;

        Subject() {}
        public Subject(int value) { alpha = value; }

        public int plus(int value) { return alpha + value; }
        private void hidden() {}
    }

    public static void main(String[] args) throws Exception {
        Class<?> cls = Subject.class;
        Field alpha = cls.getDeclaredField("alpha");
        Field beta = cls.getDeclaredField("beta");
        Method plus = cls.getDeclaredMethod("plus", int.class);
        Method hidden = cls.getDeclaredMethod("hidden");
        Constructor<?> zero = cls.getDeclaredConstructor();
        Constructor<?> one = cls.getDeclaredConstructor(int.class);
        System.out.println(cls.getDeclaredFields().length);
        System.out.println(cls.getDeclaredMethods().length);
        System.out.println(cls.getDeclaredConstructors().length);
        System.out.println(alpha.getType() == int.class);
        System.out.println(beta.getType() == String.class);
        System.out.println(plus.getReturnType() == int.class);
        System.out.println(hidden.getParameterCount());
        System.out.println(zero.getParameterCount());
        System.out.println(one.getParameterTypes()[0] == int.class);
    }
}
