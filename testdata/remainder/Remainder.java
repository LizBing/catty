public final class Remainder {
    static int rem(int x, int y) {
        return x % y;
    }

    public static void main(String[] args) {
        int x = 7;
        int y = 3;
        if (args.length == 1) x = -7;
        if (args.length == 2) y = 0;
        if (args.length == 3) y = -3;
        if (args.length == 4) {
            x = -2147483648;
            y = -1;
        }
        System.out.println(rem(x, y));
    }
}
