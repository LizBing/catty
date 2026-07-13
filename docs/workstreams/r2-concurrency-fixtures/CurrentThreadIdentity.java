public class CurrentThreadIdentity {
    public static void main(String[] args) {
        Thread first = Thread.currentThread();
        Thread second = Thread.currentThread();
        System.out.println(first == second);
        System.out.println(first.isAlive());
    }
}
