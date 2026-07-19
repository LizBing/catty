import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

class ReflectiveForeign {
    private void hidden() {}
}

public class ReflectiveFailures {
    static class Subject {
        public void explode() {
            throw new IllegalStateException("boom");
        }
    }

    public static void main(String[] args) throws Exception {
        Method explode = Subject.class.getMethod("explode");
        try {
            explode.invoke(new Object());
        } catch (IllegalArgumentException expected) {
            System.out.println("receiver");
        }
        try {
            explode.invoke(new Subject());
        } catch (InvocationTargetException expected) {
            System.out.println(expected.getCause().getClass().getName());
            System.out.println(expected.getCause().getMessage());
        }
        Method hidden = ReflectiveForeign.class.getDeclaredMethod("hidden");
        try {
            hidden.invoke(new ReflectiveForeign());
        } catch (IllegalAccessException expected) {
            System.out.println("access");
        }
    }
}
