// Category: String UTF-16 semantics — lone surrogate preservation (ADR-0023).
// Java String must preserve arbitrary code-unit sequences, including an
// unpaired surrogate. Constructed from a char[] so no literal decoding hides it.
// Expected (Temurin 25): length=1, charAt(0)=55357.
// Expected (catty R1): cannot round-trip a lone surrogate through a Go string;
// records the correctness gap. (May also be Not implemented if the char[]
// constructor is absent.)
public class LoneSurrogate {
    public static void main(String[] args) {
        char[] c = { 0xD83D };
        String s = new String(c);
        System.out.println(s.length());
        System.out.println((int) s.charAt(0));
    }
}
