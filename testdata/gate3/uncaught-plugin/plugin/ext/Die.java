package ext;

import app.DieStep;

public class Die implements DieStep {
    public static Die create() { return new Die(); }

    @Override
    public int die() { throw new RuntimeException("plugin failure"); }
}
