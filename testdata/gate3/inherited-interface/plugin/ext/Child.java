package ext;

import api.Base;
import api.Step;

public class Child extends Base implements Step {
    public static Step create() { return new Child(); }
}
