package demo;

public class Child extends Base {
    public static Child create() { return new Child(); }

    @Override
    public int value() { return 22; }
}
