import java.lang.reflect.Modifier;

public class ClassQueries {
    interface Marker {}
    static class Base {}
    static final class Child extends Base implements Marker {}

    public static void main(String[] args) {
        Class<?> cls = Child.class;
        System.out.println(cls.getSuperclass() == Base.class);
        System.out.println(cls.getInterfaces().length);
        System.out.println(cls.getInterfaces()[0] == Marker.class);
        System.out.println(Modifier.isStatic(cls.getModifiers()));
        System.out.println(Modifier.isFinal(cls.getModifiers()));
        System.out.println(Marker.class.isInterface());
    }
}
