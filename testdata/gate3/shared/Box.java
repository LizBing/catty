package app;

public class Box {
    private int value;
    private int callbacks;

    public Box(int value) { this.value = value; }

    public int add(int delta) { value += delta; callbacks++; return value; }

    public int get() { return value; }

    public int callbacks() { return callbacks; }

    public void set(int v) { value = v; }
}
