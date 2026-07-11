public class StaticFields {
    static int x = compute();
    static int y = 5;

    static int compute() {
        return 10;
    }

    public static void main(String[] args) {
        System.out.println(x);
        System.out.println(y);
        System.out.println(x + y);
    }
}
