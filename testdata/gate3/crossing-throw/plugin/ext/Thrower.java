package ext;

import app.Base;
import app.Box;
import app.Bridge;
import app.ThrowStep;

public class Thrower extends Base implements ThrowStep {
    public static Thrower create() { return new Thrower(); }

    @Override
    public int boom(Box box, Throwable thrown) throws Throwable {
        Bridge.fire(box, thrown);
        return 0;
    }
}
