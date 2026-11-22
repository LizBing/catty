package ext;

import app.Box;

public class Alpha {
    public static Alpha create() { return new Alpha(); }

    public int value(Box box) { return box.add(30); }
}
