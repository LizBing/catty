public final class ProfileImage {
    static final long LONG_CONSTANT = 0x1122334455667788L;
    static final int STATIC_FIELD = 29;

    static int mix(int value) {
        return value - 7;
    }

    static int helper(int value) {
        return value * 3 + STATIC_FIELD;
    }

    int scale(int value) {
        return value / 4 - 1;
    }

    public static void main(String[] args) {
        ProfileImage image = new ProfileImage();
        long big = LONG_CONSTANT;
        int high = (int) (big >>> 32);
        int low = (int) big;
        int folded = mix(high / 1000) + mix(low / 100000);
        int first = helper(folded);
        int scaled = image.scale(first);
        System.out.println(high);
        System.out.println(folded);
        System.out.println(first);
        System.out.println(scaled);
        if (scaled > folded) {
            System.out.println(7);
        } else {
            System.out.println(8);
        }
    }
}
