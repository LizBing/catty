package app;

public class ProbeBad extends Base {
    public static ProbeBad create() { return new ProbeBad(); }

    public int call(Base other) { return other.secret(); }
}
