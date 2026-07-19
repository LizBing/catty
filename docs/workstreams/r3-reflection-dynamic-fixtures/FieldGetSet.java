import java.lang.reflect.Field;

public class FieldGetSet {
    static class Subject {
        public int number;
        public Object reference;
        public static long total;
    }

    public static void main(String[] args) throws Exception {
        Subject subject = new Subject();
        Field number = Subject.class.getField("number");
        Field reference = Subject.class.getField("reference");
        Field total = Subject.class.getField("total");
        number.set(subject, Byte.valueOf((byte) 9));
        reference.set(subject, "value");
        total.set(null, Integer.valueOf(33));
        System.out.println(number.getInt(subject));
        System.out.println(reference.get(subject));
        System.out.println(total.getLong(null));
    }
}
