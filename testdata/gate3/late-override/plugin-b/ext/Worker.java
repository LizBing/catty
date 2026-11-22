package ext;

import app.Base;
import app.Box;

public class Worker extends Base {
    public static Worker create() { return new Worker(); }

    @Override
    public int value(Box box) { return box.add(20); }
}
