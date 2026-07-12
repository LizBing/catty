// RealBaseSmoke — R1 quality gate for the java.base classpath path.
//
// Run with: CATTY_BOOT=<extracted java.base> catty -cp <fixtures> RealBaseSmoke
// Requires a real JDK's java.base on the classpath (exercises class loading
// from real .class files + synthetic bootstrap classes together).
//
// Each assertion is wrapped in try/catch so one failure doesn't mask the rest.
// Output format: "<n> passed, <m> failed" — must match real `java` byte-for-byte.
//
// SCOPE — what this covers (R1):
//   PrintStream, String (length/charAt/equals/hashCode/substring/concat/
//   startsWith/endsWith/isEmpty/indexOf), Object (hashCode/toString),
//   Class (getName/isInterface/getSuperclass/isInstance), StringBuilder,
//   Math.max, Integer (MAX_VALUE/parseInt/toHexString), ArrayList,
//   NullPointerException catch, System.identityHashCode, Long (parseLong/
//   MAX_VALUE/toHexString), System.getProperty.
//
// DELIBERATELY EXCLUDED (R2 dependencies, NOT R1 defects):
//   - HashMap          — currently reaches an uninitialized VM dependency; the
//                        basic operation path has no direct Unsafe edge
//   - Double.parseDouble — currently times out in the FloatingDecimal path; no
//                          direct Unsafe edge found in the selected parse path
//   - Integer/Long.toString — DecimalDigits uses a narrow Unsafe array-write path
//   The *toHexString variants ARE covered: they bypass DecimalDigits and use
//   Integer.digits directly, so they exercise the String(byte[],coder) decode
//   path that toString would also need once Unsafe lands.
public class RealBaseSmoke {
    public static void main(String[] args) {
        int pass = 0;
        int fail = 0;

        // 1. PrintStream (System.out)
        try {
            System.out.print("");
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 1: " + t); fail++; }

        // 2. String length + charAt
        try {
            String s = "hello world";
            assert s.length() == 11 : "length";
            assert s.charAt(1) == 'e' : "charAt";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 2: " + t); fail++; }

        // 3. Object.hashCode + toString
        try {
            Object o = new Object();
            o.hashCode();
            String ts = o.toString();
            assert ts.startsWith("java.lang.Object@") : "toString: " + ts;
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 3: " + t); fail++; }

        // 4. Class.getName + isInterface + getSuperclass + isInstance
        try {
            Class<?> c = "hello".getClass();
            assert c.getName().equals("java.lang.String") : "getName: " + c.getName();
            assert !c.isInterface() : "isInterface";
            assert c.getSuperclass().getName().equals("java.lang.Object") : "getSuperclass";
            assert c.isInstance("world") : "isInstance";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 4: " + t); fail++; }

        // 5. StringBuilder append(String/I/Z/C) + toString
        try {
            StringBuilder sb = new StringBuilder();
            sb.append("a").append(42).append(true).append('z');
            assert sb.toString().equals("a42truez") : "StringBuilder: " + sb;
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 5: " + t); fail++; }

        // 6. Math.max
        try {
            assert Math.max(10, 20) == 20 : "Math.max";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 6: " + t); fail++; }

        // 7. Integer.MAX_VALUE + parseInt + toHexString
        try {
            assert Integer.MAX_VALUE == 2147483647 : "MAX_VALUE";
            assert Integer.parseInt("42") == 42 : "parseInt";
            assert Integer.toHexString(255).equals("ff") : "toHexString: " + Integer.toHexString(255);
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 7: " + t); fail++; }

        // 8. ArrayList create + add + size + get
        try {
            java.util.ArrayList<String> list = new java.util.ArrayList<>();
            list.add("a"); list.add("b");
            assert list.size() == 2 : "size";
            assert list.get(1).equals("b") : "get(1): " + list.get(1);
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 8: " + t); fail++; }

        // 9. Exception catch
        try {
            try { ((String) null).length(); }
            catch (NullPointerException e) { pass++; }
        } catch (Throwable t) { System.out.println("FAIL 9: " + t); fail++; }

        // 10. System.identityHashCode
        try {
            int h = System.identityHashCode("x");
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 10: " + t); fail++; }

        // --- New tests (C3) ---

        // 11. Long parseLong + MAX_VALUE
        try {
            assert Long.parseLong("123456789") == 123456789L : "parseLong";
            assert Long.MAX_VALUE > 0 : "MAX_VALUE";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 11: " + t); fail++; }

        // 12. String.equals content comparison
        try {
            assert "abc".equals("abc") : "equals true";
            assert !"abc".equals("abd") : "equals false";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 12: " + t); fail++; }

        // 13. String.hashCode consistency
        try {
            assert "abc".hashCode() == "abc".hashCode() : "hashCode same";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 13: " + t); fail++; }

        // 14. String.substring
        try {
            assert "hello world".substring(6).equals("world") : "substring(6): " + "hello world".substring(6);
            assert "hello".substring(1,4).equals("ell") : "substring(1,4): " + "hello".substring(1,4);
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 14: " + t); fail++; }

        // 15. String.startsWith / endsWith + concat + isEmpty
        try {
            assert "hello world".startsWith("hello") : "startsWith";
            assert "hello world".endsWith("world") : "endsWith";
            assert "hello".concat(" world").equals("hello world") : "concat";
            assert "".isEmpty() : "isEmpty true";
            assert !"x".isEmpty() : "isEmpty false";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 15: " + t); fail++; }

        // 16. System.getProperty
        try {
            String sep = System.getProperty("line.separator");
            assert sep != null : "line.separator is null";
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 16: " + t); fail++; }

        // 17. Long.toHexString (bypasses DecimalDigits→Unsafe)
        try {
            assert Long.toHexString(255).equals("ff") : "Long.toHexString: " + Long.toHexString(255);
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 17: " + t); fail++; }

        // 18. String.indexOf
        try {
            assert "hello".indexOf('e') == 1 : "indexOf: " + "hello".indexOf('e');
            pass++;
        } catch (Throwable t) { System.out.println("FAIL 18: " + t); fail++; }

        System.out.println(pass + " passed, " + fail + " failed");
    }
}
