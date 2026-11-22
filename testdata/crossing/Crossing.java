import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public final class Crossing {
    public interface Step {
        void apply(Box box);
    }

    public static final class Signal extends RuntimeException {}

    public static final class Box {
        int count;
        final Signal signal = new Signal();

        public void bump() {
            count++;
        }
    }

    public static final class Worker implements Step {
        public void apply(Box box) {
            box.bump();
            throw box.signal;
        }
    }

    public static int probe() throws ReflectiveOperationException {
        Box box = new Box();
        Step step = new Worker();
        try {
            step.apply(box);
            return -1;
        } catch (Signal e) {
            if (e != box.signal || box.count != 1) return -2;
        }
        Method method = Step.class.getMethod("apply", Box.class);
        try {
            method.invoke(step, box);
            return -3;
        } catch (InvocationTargetException e) {
            if (e.getCause() != box.signal) return -4;
        }
        return box.count;
    }
}
