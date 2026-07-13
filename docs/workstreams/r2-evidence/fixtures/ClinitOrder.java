// Category: class initialization ordering (JVMS §5.5).
// Probes: superclass <clinit> runs before subclass <clinit>.
// Expected (Temurin 25 and catty): "A" then "B".
class A {
    static { System.out.println("A"); }
}
class B extends A {
    static { System.out.println("B"); }
}
public class ClinitOrder {
    public static void main(String[] args) {
        new B();
        System.out.println("main");
    }
}
