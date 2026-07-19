public class StatelessLambda {
    interface IntOperation {
        int apply(int value);
    }

    static IntOperation operation() {
        return value -> value * value;
    }

    public static void main(String[] args) {
        IntOperation first = operation();
        IntOperation second = operation();
        System.out.println(first.apply(7));
        System.out.println(second.apply(-3));
        System.out.println(first.getClass() == second.getClass());
    }
}
