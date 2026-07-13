// Category: class initialization trigger — invokestatic (JVMS §5.5).
// Probes: invokestatic triggers <clinit> before the method runs.
// Expected (Temurin 25 and catty): "Holder" then "42".
class Holder {
    static int x;
    static {
        System.out.println("Holder");
        x = 42;
    }
    static int getX() {
        return x;
    }
}
public class InvokeStaticInit {
    public static void main(String[] args) {
        System.out.println(Holder.getX());
    }
}
