// Category: String UTF-16 semantics — bounds failures.
// Expected (Temurin 25): StringIndexOutOfBoundsException twice.
public class StringBounds {
    public static void main(String[] args) {
        String value = "x";
        try {
            value.charAt(1);
        } catch (Throwable failure) {
            System.out.println(failure.getClass().getName());
        }
        try {
            value.substring(1, 0);
        } catch (Throwable failure) {
            System.out.println(failure.getClass().getName());
        }
    }
}
