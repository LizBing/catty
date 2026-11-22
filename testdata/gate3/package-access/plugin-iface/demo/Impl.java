package demo;

public class Impl implements HiddenIface {
    public static Impl create() { return new Impl(); }

    @Override
    public int value() { return 2; }
}
