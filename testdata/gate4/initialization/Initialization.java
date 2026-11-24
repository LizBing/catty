public class Initialization {
    private static final RuntimeException TOKEN = new RuntimeException();
    private static volatile boolean initializing;
    private static volatile boolean attempting;
    private static volatile boolean release;
    private static int attempts;
    private static int ownerResult;
    private static int contenderResult;

    static class Ready {
        static int value = initialize();

        private static int initialize() {
            attempts++;
            Check.that(read() == 0, "recursive initialization should see the default");
            blockInitializer();
            return 73;
        }

        static int read() {
            return value;
        }
    }

    static class Broken {
        static int value = initialize();

        private static int initialize() {
            attempts++;
            blockInitializer();
            throw TOKEN;
        }

        static int read() {
            return value;
        }
    }

    static class Reader extends Thread {
        private final boolean owner;
        private final boolean broken;

        Reader(boolean owner, boolean broken) {
            this.owner = owner;
            this.broken = broken;
        }

        public void run() {
            if (!owner) {
                Thread.currentThread().interrupt();
                attempting = true;
            }
            int result;
            try {
                result = broken ? Broken.read() : Ready.read();
            } catch (ExceptionInInitializerError error) {
                Check.that(error.getCause() == TOKEN, "initializer lost original cause");
                result = -1;
            } catch (NoClassDefFoundError error) {
                result = -2;
            }
            if (owner) {
                ownerResult = result;
            } else {
                Check.that(Thread.interrupted(), "initialization consumed interrupt");
                contenderResult = result;
            }
        }
    }

    private static void blockInitializer() {
        initializing = true;
        while (!release) {
            Thread.yield();
        }
    }

    private static void trial(boolean broken) throws InterruptedException {
        initializing = false;
        attempting = false;
        release = false;
        Reader owner = new Reader(true, broken);
        Reader contender = new Reader(false, broken);
        owner.start();
        while (!initializing) {
            Thread.yield();
        }
        contender.start();
        while (!attempting) {
            Thread.yield();
        }
        release = true;
        owner.join();
        contender.join();
        Check.that(ownerResult == (broken ? -1 : 73), "bad initializer owning outcome");
        Check.that(contenderResult == (broken ? -2 : 73), "bad initializer competing outcome");
    }

    public static void main(String[] args) throws InterruptedException {
        Class<?> readyType = Ready.class;
        Class<?> brokenType = Broken.class;
        Check.that(readyType != brokenType && attempts == 0, "class literal triggered initialization");
        trial(false);
        trial(true);
        try {
            Broken.read();
            throw new RuntimeException("erroneous class ran again");
        } catch (NoClassDefFoundError expected) {
            Check.that(attempts == 2, "initializer did not run exactly once per class");
        }
        System.out.println("initialization");
    }
}
