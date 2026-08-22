package app;

public class Main {
    public static void main(String[] args) {
        Greeter g = new Helper();
        System.out.println(g.greet("m1"));
        System.out.println(g.greet("m2"));
        Object o = new Helper();
        System.out.println(o instanceof Greeter);
        Helper h = (Helper) o;
        System.out.println(h.hits());
    }
}
