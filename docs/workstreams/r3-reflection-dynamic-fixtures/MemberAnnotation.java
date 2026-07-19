import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;

public class MemberAnnotation {
    @Retention(RetentionPolicy.RUNTIME)
    @Target({ElementType.FIELD, ElementType.METHOD, ElementType.PARAMETER})
    @interface Mark {
        int value();
    }

    static class Subject {
        @Mark(1)
        public int field;

        @Mark(2)
        public void method(@Mark(3) String input) {}
    }

    public static void main(String[] args) throws Exception {
        Field field = Subject.class.getField("field");
        Method method = Subject.class.getMethod("method", String.class);
        Parameter parameter = method.getParameters()[0];
        System.out.println(field.getAnnotation(Mark.class).value());
        System.out.println(method.getAnnotation(Mark.class).value());
        System.out.println(parameter.getAnnotation(Mark.class).value());
    }
}
