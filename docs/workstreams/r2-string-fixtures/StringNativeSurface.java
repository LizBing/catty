// Category: String native surface — null contract, PrintStream, StringBuilder.
//
// Groups:
//   A: NPE on null constructor / method args (7 cases: 6 String + 1 System)
//   B: "null" output for PrintStream.println/print(String null) and
//      StringBuilder.append(String null)
//   C: Lone surrogate println(char) -> '?'
//
// Expected (Temurin 25): see per-group comments below.
//
// NOTE: No string concatenation (+) is used in catch handlers — that would
// generate invokedynamic (StringConcatFactory) which is not yet supported.
public class StringNativeSurface {
    // ---- Group A: NPE ----
    static void testNPE() {
        // new String((String) null)
        try { new String((String) null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // new String((char[]) null)
        try { new String((char[]) null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // concat(null)
        try { "x".concat(null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // startsWith(null)
        try { "x".startsWith(null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // endsWith(null)
        try { "x".endsWith(null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // compareTo(null)
        try { "x".compareTo(null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
        // System.getProperty(null)
        try { System.getProperty(null); System.out.println("A-FAIL"); } catch (Throwable t) { System.out.println("A:NPE"); }
    }

    // ---- Group B: "null" output ----
    static void testNullOutput() {
        System.out.println((String) null);
        System.out.print((String) null);
        System.out.println();
        StringBuilder sb = new StringBuilder();
        sb.append((String) null);
        sb.append("|");
        sb.append("X");
        System.out.println(sb.toString());
    }

    // ---- Group C: lone surrogate println(char) ----
    static void testLoneSurrogateChar() {
        System.out.print("lone-high:");
        System.out.println((char) 0xD800);
        System.out.print("lone-low:");
        System.out.println((char) 0xDC00);
        System.out.print("valid:");
        System.out.println('A');
    }

    // ---- main ----
    public static void main(String[] args) {
        testNPE();
        testNullOutput();
        testLoneSurrogateChar();
    }
}
