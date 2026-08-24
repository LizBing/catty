// DEBT-0017 probe: a static native returning category-2 (long) must flow
// through the emitted call site with correct operand-stack depth. The
// historical limitation forced Bench to use an int32 clock workaround.
public class JProbe {
    public static void main(String[] args) {
        long t = System.nanoTime();
        long d = t - t + 42;
        System.out.println(d);
    }
}
