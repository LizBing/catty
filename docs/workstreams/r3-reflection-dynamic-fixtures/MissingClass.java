public class MissingClass {
    public static void main(String[] args) {
        try {
            Class.forName("catty.missing.NeverDefined");
            System.out.println("unexpected");
        } catch (ClassNotFoundException expected) {
            System.out.println(expected.getClass().getName());
        }
    }
}
