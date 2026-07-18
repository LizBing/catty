import java.lang.annotation.Inherited;
import java.lang.annotation.Repeatable;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

public class InheritedRepeatableAnnotation {
    @Inherited
    @Retention(RetentionPolicy.RUNTIME)
    @Repeatable(Tags.class)
    @interface Tag {
        String value();
    }

    @Inherited
    @Retention(RetentionPolicy.RUNTIME)
    @interface Tags {
        Tag[] value();
    }

    @Tag("first")
    @Tag("second")
    static class Base {}

    static class Child extends Base {}

    public static void main(String[] args) {
        Tag[] tags = Child.class.getAnnotationsByType(Tag.class);
        System.out.println(tags.length);
        System.out.println(tags[0].value());
        System.out.println(tags[1].value());
        System.out.println(Child.class.getDeclaredAnnotationsByType(Tag.class).length);
    }
}
