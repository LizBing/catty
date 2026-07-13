// Category: String UTF-16 semantics — char[] constructor and toCharArray copy.
// Expected (Temurin 25): 3, 55357, 56832, 55357, 55357.
public class StringCharArrayRoundTrip {
    public static void main(String[] args) {
        char[] input = { 0xD83D, 0xDE00, 0xD83D };
        String value = new String(input);
        input[0] = 'X';
        char[] copy = value.toCharArray();
        copy[2] = 'Y';
        System.out.println(value.length());
        System.out.println((int) value.charAt(0));
        System.out.println((int) value.charAt(1));
        System.out.println((int) value.charAt(2));
        System.out.println((int) value.toCharArray()[0]);
    }
}
