package app;

public class Helper implements Greeter {
    private int hits;

    public String greet(String who) {
        synchronized (this) {
            synchronized (this) {
                hits++;
            }
        }
        return "hi:" + who + ":" + hits;
    }

    public int hits() {
        return hits;
    }
}
