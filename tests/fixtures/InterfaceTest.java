public class InterfaceTest {
    public static void main(String[] args) {
        // 1. Interface dispatch via Comparable
        Shape[] shapes = new Shape[3];
        shapes[0] = new Shape(5);
        shapes[1] = new Shape(2);
        shapes[2] = new Shape(8);

        // Sort by area (uses invokeinterface for Comparable.compareTo)
        for (int i = 0; i < shapes.length; i++) {
            for (int j = i + 1; j < shapes.length; j++) {
                if (shapes[i].compareTo(shapes[j]) > 0) {
                    Shape tmp = shapes[i];
                    shapes[i] = shapes[j];
                    shapes[j] = tmp;
                }
            }
        }

        for (int i = 0; i < shapes.length; i++) {
            System.out.println("area: " + shapes[i].area());
        }

        // 2. Multi-dimensional array
        int[][] grid = new int[2][3];
        for (int i = 0; i < 2; i++) {
            for (int j = 0; j < 3; j++) {
                grid[i][j] = i * 3 + j;
            }
        }
        for (int i = 0; i < 2; i++) {
            for (int j = 0; j < 3; j++) {
                System.out.println("grid[" + i + "][" + j + "] = " + grid[i][j]);
            }
        }

        System.out.println("done");
    }
}

class Shape implements Comparable<Shape> {
    int size;

    Shape(int size) {
        this.size = size;
    }

    int area() {
        return size * size;
    }

    public int compareTo(Shape other) {
        return this.area() - other.area();
    }
}
