package demo;

public class LocalImpl implements LocalIface {
    public static LocalImpl create() { return new LocalImpl(); }

    @Override
    public int value() { return 4; }
}
