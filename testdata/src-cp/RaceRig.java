// P-0009 T5 experiment: two threads racing an unsynchronized int field.
// Built with `go build -race`, the Go race detector should observe the
// concurrent unsynchronized accesses through the shared Instance fields
// slice — "白嫖" observability claim, selling point #5.
public class RaceRig {
    static int counter = 0;

    static class Bumper extends Thread {
        public void run() {
            for (int i = 0; i < 100000; i++) {
                counter++;
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Bumper a = new Bumper();
        Bumper b = new Bumper();
        a.start();
        b.start();
        a.join();
        b.join();
        System.out.println("counter=" + counter);
    }
}
