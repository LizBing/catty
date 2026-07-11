public class ExceptionTest {
    public static void main(String[] args) {
        // 1. Explicit throw + catch
        try {
            throw new RuntimeException("explicit throw");
        } catch (RuntimeException e) {
            System.out.println("caught: " + e.getMessage());
        }

        // 2. NPE
        try {
            String s = null;
            int len = s.length();
            System.out.println("should not print: " + len);
        } catch (NullPointerException e) {
            System.out.println("caught NPE");
        }

        // 3. ArithmeticException (division by zero)
        try {
            int x = 10;
            int y = 0;
            int z = x / y;
            System.out.println("should not print: " + z);
        } catch (ArithmeticException e) {
            System.out.println("caught: " + e.getMessage());
        }

        // 4. Finally always runs
        try {
            System.out.println("in try");
        } finally {
            System.out.println("in finally");
        }

        // 5. Finally runs even when exception thrown
        try {
            throw new IllegalArgumentException("boom");
        } catch (IllegalArgumentException e) {
            System.out.println("caught: " + e.getMessage());
        } finally {
            System.out.println("finally after catch");
        }

        // 6. Exception propagates through method call
        try {
            risky();
            System.out.println("should not print");
        } catch (ArithmeticException e) {
            System.out.println("caught from risky(): " + e.getMessage());
        }

        System.out.println("done");
    }

    static void risky() {
        int x = 10;
        int y = 0;
        int z = x / y;
        System.out.println("should not print: " + z);
    }
}
