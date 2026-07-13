public class MonitorOwnership {
    public static void main(String[] args) {
        Object lock = new Object();
        try {
            lock.notify();
            System.out.println("missing");
        } catch (IllegalMonitorStateException expected) {
            System.out.println("IllegalMonitorStateException");
        }
    }
}
