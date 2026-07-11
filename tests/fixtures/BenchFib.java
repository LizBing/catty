public class BenchFib {
    // fib(35) = 9227465; ~29 million recursive calls. Pure interpreter work,
    // so it stresses the dispatch loop, method entry/exit, and int arithmetic.
    static int fib(int n) {
        if (n < 2) return n;
        return fib(n - 1) + fib(n - 2);
    }

    public static void main(String[] args) {
        System.out.println(fib(35));
    }
}
