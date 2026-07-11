public class Factorial {
    static long fact(long n) {
        if (n <= 1) return 1;
        return n * fact(n - 1);
    }

    public static void main(String[] args) {
        long[] results = new long[10];
        for (int i = 0; i < 10; i++) {
            results[i] = fact(i);
        }
        for (int i = 0; i < results.length; i++) {
            System.out.println(i + "! = " + results[i]);
        }
    }
}
