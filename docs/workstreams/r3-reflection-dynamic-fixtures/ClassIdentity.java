public class ClassIdentity {
    public static void main(String[] args) throws Exception {
        Class<?> literal = String.class;
        Class<?> loaded = Class.forName("java.lang.String");
        String value = "catty";
        System.out.println(literal == loaded);
        System.out.println(literal == value.getClass());
        System.out.println(literal.getName());
    }
}
