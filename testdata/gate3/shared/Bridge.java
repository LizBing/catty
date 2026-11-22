package app;

public class Bridge {
    public static void fire(Box box, Throwable thrown) throws Throwable {
        box.add(100);
        throw thrown;
    }
}
