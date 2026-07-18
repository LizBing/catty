public class MethodReference {
    interface StringOperation {
        int apply(String value);
    }

    interface BoundOperation {
        String apply(String value);
    }

    static String decorate(String prefix, String value) {
        return prefix + value;
    }

    public static void main(String[] args) {
        StringOperation length = String::length;
        String prefix = "catty-";
        BoundOperation decorate = value -> decorate(prefix, value);
        System.out.println(length.apply("reflection"));
        System.out.println(decorate.apply("dynamic"));
    }
}
