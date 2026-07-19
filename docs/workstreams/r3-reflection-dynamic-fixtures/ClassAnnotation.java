import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

public class ClassAnnotation {
    @Retention(RetentionPolicy.RUNTIME)
    @interface Label {
        String name();
        int version();
    }

    @Label(name = "catty", version = 3)
    static class Subject {}

    public static void main(String[] args) {
        Label label = Subject.class.getAnnotation(Label.class);
        System.out.println(label != null);
        System.out.println(label.name());
        System.out.println(label.version());
        System.out.println(Subject.class.isAnnotationPresent(Label.class));
    }
}
