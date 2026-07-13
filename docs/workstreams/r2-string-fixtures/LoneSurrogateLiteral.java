// Category: String UTF-16 semantics — lone surrogate in classfile literal.
// The Unicode escape \uD800 is processed at the lexical level (JLS §3.3) and
// becomes a U+D800 code point in the string literal. javac encodes it as
// 0xED 0xA0 0x80 in the CONSTANT_Utf8 constant pool (3-byte MUTF-8).
// decodeMUTF8ToUTF16 must preserve it as a single 0xD800 code unit.
//
// Expected (Temurin 25): 1, 55296.
public class LoneSurrogateLiteral {
    public static void main(String[] args) {
        String s = "\uD800"; // lone high surrogate from a literal
        System.out.println(s.length());
        System.out.println((int) s.charAt(0));
    }
}
