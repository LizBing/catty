package ext;

import app.Base;
import app.Box;

public class Worker extends Base {
    static int counter;

    public static Worker create() { return new Worker(); }

    @Override
    public int value(Box box) { counter += 20; return counter; }
}
