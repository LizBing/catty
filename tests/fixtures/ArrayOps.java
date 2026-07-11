public class ArrayOps {
    // Merge-free: aload, iconst, iaload, ireturn — exercises a ref slot (the
    // array) and an int slot, which under fresh-per-def become distinct typed
    // Go temps (*rtda.Object / int32).
    static int first(int[] a) {
        return a[0];
    }

    // Has a merge (loop back-edge): used to exercise the merge-free gate, which
    // should refuse to AOT this method in A2.1 (phis are A2.3).
    static int sum(int[] a) {
        int s = 0;
        for (int i = 0; i < a.length; i++) {
            s += a[i];
        }
        return s;
    }

    // A diamond: `a > b ? a : b` leaves a value on the operand stack across the
    // join, so it needs phi — used to exercise the non-empty-stack-merge gate.
    static int max(int a, int b) {
        return a > b ? a : b;
    }

    // Float (category-1, float32) — straight-line, no merges.
    static float fadd(float a, float b) {
        return a + b;
    }

    // Double (category-2, float64) — straight-line, no merges.
    static double dmul(double a, double b) {
        return a * b;
    }

    // Float remainder — Go has no float %, uses runtime.FloatMod (math.Mod).
    static float frem(float a, float b) {
        return a % b;
    }

    // Long crossing a diamond — exercises cat-2 merge phi (int64 at the join).
    static long lcond(boolean c, long a, long b) {
        return c ? a : b;
    }

    // Switch (tableswitch) — dense keys [1..2] + default.
    static int sw(int n) {
        switch (n) {
            case 1: return 10;
            case 2: return 20;
            default: return 0;
        }
    }
}
