package demo;

public class PublicChild extends PublicBase {
    public static PublicChild create() { return new PublicChild(); }

    @Override
    public int value() { return 22; }
}
