// Category: class initialization — ConstantValue is not an init trigger.
// Expected (Temurin 25): "7" only; Constants.<clinit> does not run.
class Constants {
    static final int VALUE = 7;
    static int marker = mark();
    static int mark() {
        System.out.println("Constants");
        return 1;
    }
}
public class ConstantFieldNoInit {
    public static void main(String[] args) {
        System.out.println(Constants.VALUE);
    }
}
