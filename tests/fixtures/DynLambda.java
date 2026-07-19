// JDK 25 integration fixture: capturing lambda produces invokedynamic.
import java.util.function.IntUnaryOperator;

public class DynLambda {
    public static IntUnaryOperator makeAdder(int base) {
        return n -> base + n;
    }

    public static void main(String[] args) {
        IntUnaryOperator add5 = makeAdder(5);
        System.out.println(add5.applyAsInt(3));
    }
}
