// Category: String UTF-16 semantics — hashCode (ADR-0023).
// Java hashCode is defined over UTF-16 code units: h = 31*h + unit.
// For U+1F600 (units 0xD83D, 0xDE00) the Temurin-25 result is shown below;
// any rune-based or byte-based hash diverges.
// Expected (Temurin 25): 1772899  (computed: 31*0xD83D + 0xDE00)
// Expected (catty R1): diverges (rune/byte-based hash).
public class HashDivergence {
    public static void main(String[] args) {
        String s = "😀";
        System.out.println(s.hashCode());
    }
}
