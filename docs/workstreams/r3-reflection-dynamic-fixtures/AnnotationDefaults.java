import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

public class AnnotationDefaults {
    enum Level { LOW, HIGH }

    @Retention(RetentionPolicy.RUNTIME)
    @interface Nested {
        String value() default "nested";
    }

    @Retention(RetentionPolicy.RUNTIME)
    @interface Config {
        int number() default 7;
        Level level() default Level.HIGH;
        Class<?> type() default String.class;
        Nested nested() default @Nested;
        int[] values() default {2, 4, 6};
    }

    @Config
    static class Subject {}

    public static void main(String[] args) {
        Config config = Subject.class.getAnnotation(Config.class);
        System.out.println(config.number());
        System.out.println(config.level().name());
        System.out.println(config.type() == String.class);
        System.out.println(config.nested().value());
        for (int value : config.values()) {
            System.out.println(value);
        }
    }
}
