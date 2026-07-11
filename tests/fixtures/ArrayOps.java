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
}
