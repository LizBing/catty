public class InterruptStatus {
    public static void main(String[] args) {
        Thread current = Thread.currentThread();
        current.interrupt();
        System.out.println(current.isInterrupted());
        System.out.println(Thread.interrupted());
        System.out.println(current.isInterrupted());
    }
}
