public class CapturingLambda {
    interface LongOperation {
        long apply(long value);
    }

    static LongOperation add(long captured) {
        return value -> captured + value;
    }

    public static void main(String[] args) {
        LongOperation ten = add(10L);
        LongOperation negative = add(-5L);
        System.out.println(ten.apply(7L));
        System.out.println(negative.apply(7L));
        System.out.println(ten.getClass() == negative.getClass());
        System.out.println(ten != negative);
    }
}
