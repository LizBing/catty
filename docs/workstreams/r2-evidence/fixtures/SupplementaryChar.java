// Category: String UTF-16 semantics — supplementary character (ADR-0023).
// U+1F600 is encoded as two UTF-16 code units: 0xD83D (55357) + 0xDE00 (56832).
// Expected (Temurin 25): length=2, charAt(0)=55357, charAt(1)=56832.
// Expected (catty R1): diverges — length/charAt are not UTF-16-code-unit based.
public class SupplementaryChar {
    public static void main(String[] args) {
        String s = "😀";   // U+1F600, written as a surrogate pair
        System.out.println(s.length());
        System.out.println((int) s.charAt(0));
        System.out.println((int) s.charAt(1));
    }
}
