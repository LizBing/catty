public class PrimitiveAndArrayClass {
    public static void main(String[] args) throws Exception {
        System.out.println(int.class.isPrimitive());
        System.out.println(void.class.isPrimitive());
        System.out.println(int[].class.isArray());
        System.out.println(int[].class.getComponentType() == int.class);
        System.out.println(String[][].class.getComponentType() == String[].class);
        System.out.println(Class.forName("[I") == int[].class);
        System.out.println(Class.forName("[[Ljava.lang.String;") == String[][].class);
    }
}
