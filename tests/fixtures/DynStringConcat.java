// JDK 25 integration fixture: string concatenation produces invokedynamic.
public class DynStringConcat {
    public static String greet(String name) {
        return "Hello, " + name + "!";
    }

    public static void main(String[] args) {
        System.out.println(greet("catty"));
    }
}
