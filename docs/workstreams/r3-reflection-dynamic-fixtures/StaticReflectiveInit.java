import java.lang.reflect.Field;
import java.lang.reflect.Method;

public class StaticReflectiveInit {
    static class FieldTarget {
        static int value = initialize();

        static int initialize() {
            System.out.println("field-init");
            return 17;
        }
    }

    static class MethodTarget {
        static int marker = initialize();

        static int initialize() {
            System.out.println("method-init");
            return 0;
        }

        public static int answer() {
            return 29;
        }
    }

    public static void main(String[] args) throws Exception {
        Field field = FieldTarget.class.getDeclaredField("value");
        Method method = MethodTarget.class.getDeclaredMethod("answer");
        System.out.println("metadata-ready");
        System.out.println(field.getInt(null));
        System.out.println(method.invoke(null));
    }
}
