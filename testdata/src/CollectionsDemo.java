import java.util.ArrayList;

public class CollectionsDemo {
    public static void main(String[] args) {
        ArrayList<Integer> xs = new ArrayList<Integer>();
        for (int i = 1; i <= 5; i++) {
            xs.add(i * 10);
        }
        int sum = 0;
        for (int i = 0; i < xs.size(); i++) {
            sum += xs.get(i);
        }
        System.out.println(sum);

        try {
            System.out.println(xs.get(99));
        } catch (IndexOutOfBoundsException e) {
            System.out.println("caught: " + e.getMessage());
        }

        try {
            int z = sum / (5 - 5);
            System.out.println(z);
        } catch (ArithmeticException e) {
            System.out.println("div-by-zero caught");
        }

        StringBuilder sb = new StringBuilder();
        sb.append("sum=").append(sum);
        System.out.println(sb.toString());
        System.out.println(xs.contains(30) ? "has30" : "no30");
    }
}
