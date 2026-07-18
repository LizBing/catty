import java.lang.reflect.Constructor;

public class ConstructorInvoke {
    static class Subject {
        final String name;
        final int value;

        public Subject(String name, int value) {
            this.name = name;
            this.value = value;
        }
    }

    public static void main(String[] args) throws Exception {
        Constructor<Subject> ctor = Subject.class.getConstructor(String.class, int.class);
        Subject subject = ctor.newInstance("catty", Integer.valueOf(25));
        System.out.println(subject.getClass() == Subject.class);
        System.out.println(subject.name);
        System.out.println(subject.value);
    }
}
