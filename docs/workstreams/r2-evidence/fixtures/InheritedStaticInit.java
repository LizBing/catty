// Category: class initialization trigger — inherited invokestatic (JVMS §5.5).
// Probes: invokestatic on a subclass must initialize the actual declaring
// class of the resolved method (the ancestor), not the constant-pool
// referenced class (the descendant). Init must run exactly once.
// Expected (Temurin 25 and catty): "Ancestor" then "99" then "99".
class Ancestor {
    static { System.out.println("Ancestor"); }
    static int val() { return 99; }
}
class Descendant extends Ancestor {
    // inherits val() — val() is resolved to Ancestor.val()
}
public class InheritedStaticInit {
    public static void main(String[] args) {
        System.out.println(Descendant.val());
        System.out.println(Descendant.val());
    }
}
