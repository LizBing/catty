public class SynchronizedMethods {
    synchronized boolean instanceLockHeld() {
        return Thread.holdsLock(this);
    }

    static synchronized boolean classLockHeld() {
        return Thread.holdsLock(SynchronizedMethods.class);
    }

    public static void main(String[] args) {
        System.out.println(new SynchronizedMethods().instanceLockHeld());
        System.out.println(classLockHeld());
    }
}
