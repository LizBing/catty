package demo;

public class LocalChild extends LocalParent {
    public static LocalChild create() { return new LocalChild(); }

    @Override
    public int value() { return 3; }
}
