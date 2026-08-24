// Second stack-backfill fixture: fire() throws a user exception built
// through the IllegalStateException("<init>") chain — the captured trace
// must show fire/main and must NOT contain <init> frames
// (fillInStackTrace semantics).
public class TraceProbe2 {
    public static void main(String[] args) {
        fire(3);
    }

    static void fire(int x) {
        throw new IllegalStateException("boom:" + x);
    }
}
