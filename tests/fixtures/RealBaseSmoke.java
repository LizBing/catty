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

        // 7. Integer.MAX_VALUE + parseInt
        try {
            assert Integer.MAX_VALUE == 2147483647 : "MAX_VALUE";
            assert Integer.parseInt("42") == 42 : "parseInt";
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

        System.out.println(pass + " passed, " + fail + " failed");
    }
}
