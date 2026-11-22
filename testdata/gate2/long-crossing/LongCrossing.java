public final class LongCrossing {
    static long leaf(long value, int delta) {
        return value + delta;
    }

    static long bridge(long value, int delta) {
        long saved = value ^ 0x0102030405060708L;
        if (delta < 0) {
            return leaf(saved, -delta);
        }
        return leaf(saved, delta);
    }

    public static void main(String[] args) {
        long result = bridge(0x1122334455667788L, args.length - 1);
        System.out.println((int) (result >>> 32));
        System.out.println((int) result);
    }
}
