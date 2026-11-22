package app;

public class Base implements Step {
    @Override
    public int value(Box box) {
        return box.add(1);
    }
}
