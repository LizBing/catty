// Category: String UTF-16 semantics — substring indices are code-unit indices.
// Expected (Temurin 25): 2, 55357, 56832, 1, 55357.
public class StringSubstringUnits {
    public static void main(String[] args) {
        String value = "A😀B";
        String pair = value.substring(1, 3);
        String high = value.substring(1, 2);
        System.out.println(pair.length());
        System.out.println((int) pair.charAt(0));
        System.out.println((int) pair.charAt(1));
        System.out.println(high.length());
        System.out.println((int) high.charAt(0));
    }
}
