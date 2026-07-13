public class MonitorNull {
    public static void main(String[] args) {
        Object lock = null;
        try {
            synchronized (lock) {
                System.out.println("missing");
            }
        } catch (NullPointerException expected) {
            System.out.println("NullPointerException");
        }
    }
}
