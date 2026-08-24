// Compile-time stub for the Go-registered bridge. At runtime the kernel's
// synthesized go/Bridge (registered via DefineClass) shadows this class —
// registry lookup precedes the classpath resolver.
package go;

public class Bridge {
    public static int add(int a, int b) { throw new UnsupportedOperationException("stub"); }
    public static String greet(String name) { throw new UnsupportedOperationException("stub"); }
    public static void fail() { throw new UnsupportedOperationException("stub"); }
}
