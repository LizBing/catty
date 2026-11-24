package gate5bench;

/** Public reflection target: a public class with a public static method. */
public class ReflectTarget {
    public static int combine(int left, int right) {
        return left * 31 + right;
    }
}
