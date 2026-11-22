package ext;

import app.Base;
import app.Box;

public class Child extends Base {
    public static Child create() { return new Child(); }

    @Override
    public int value(Box box) { return box.add(2); }
}
