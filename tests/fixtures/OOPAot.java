class Box {
    int v;

    int doubled() {
        return v + v;
    }
}

public class OOPAot {
    public static void main(String[] args) {
        Box b = new Box();
        b.v = 21;
        System.out.println(b.v + b.doubled());
    }
}
