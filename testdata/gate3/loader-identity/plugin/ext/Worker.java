package ext;

import app.Base;
import app.Box;

public class Worker extends Base {
    static int counter;

    public static Worker create() { return new Worker(); }

    public static boolean self(Object o) { return o instanceof Worker; }

    public static Object cast(Object o) { return (Worker) o; }

    public static int store(Object o) {
        Object[] array = new Worker[1];
        array[0] = o;
        return array.length;
    }

    @Override
    public int value(Box box) { counter++; return counter; }
}
