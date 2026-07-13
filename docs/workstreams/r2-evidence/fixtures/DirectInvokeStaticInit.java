// Category: class initialization trigger — direct AOT invokestatic (JVMS §5.5).
// Probes: invokestatic triggers <clinit> even when the called method does NOT
// read or write static state. The direct AOT call path must emit an init guard
// before the call; without it <clinit> would never run.
// Expected (Temurin 25 and catty): "SideEffect" then "1" then "2".
class SideEffect {
    static { System.out.println("SideEffect"); }
    static int identity(int x) { return x; }
}
public class DirectInvokeStaticInit {
    public static void main(String[] args) {
        System.out.println(SideEffect.identity(1));
        System.out.println(SideEffect.identity(2));
    }
}
