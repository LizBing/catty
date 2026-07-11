class Point {
    int x;
    int y;

    Point(int x, int y) {
        this.x = x;
        this.y = y;
    }

    int manhattan() {
        return x + y;
    }
}

public class OOPDemo {
    public static void main(String[] args) {
        Point p = new Point(3, 4);
        System.out.println(p.manhattan());

        Point q = new Point(10, 20);
        int total = q.x + q.manhattan();
        System.out.println(total);

        System.out.println(p.x == 3);
    }
}
