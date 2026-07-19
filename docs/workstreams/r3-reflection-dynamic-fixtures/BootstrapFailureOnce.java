public class BootstrapFailureOnce {
    interface Action {
        void run();
    }

    static Action make() {
        return BootstrapFailureTarget::run;
    }

    public static void main(String[] args) {
        Throwable first = null;
        Throwable second = null;
        try {
            make();
        } catch (Throwable failure) {
            first = failure;
            System.out.println(failure.getClass().getName());
        }
        try {
            make();
        } catch (Throwable failure) {
            second = failure;
            System.out.println(failure.getClass().getName());
        }
        System.out.println(first != null);
        System.out.println(second != null);
        System.out.println(first == second);
    }
}

class BootstrapFailureTarget {
    static void run() {}
}
