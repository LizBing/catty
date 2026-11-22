package app;

public class Base implements Step {
    public static int imageMagic() { return 41; }

    protected int secret() { return 7; }

    int tag() { return 1; }

    @Override
    public int value(Box box) { return box.add(1); }
}
