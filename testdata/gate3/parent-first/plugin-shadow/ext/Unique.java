package ext;

import app.Box;
import app.Step;

public class Unique implements Step {
    public static Unique create() { return new Unique(); }

    @Override
    public int value(Box box) { return box.add(100); }
}
