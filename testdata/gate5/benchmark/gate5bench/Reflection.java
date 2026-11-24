package gate5bench;

import java.lang.reflect.Method;

/**
 * Reflection kernels. Lookup by name is timed as its own kernel and is never
 * inside the cached-invoke window: the invoke kernel receives a Method object
 * resolved once before any timing starts.
 */
final class Reflection {
    private Reflection() {
    }

    static Method lookupCombine() throws Exception {
        return ReflectTarget.class.getMethod("combine",
                new Class[]{Integer.TYPE, Integer.TYPE});
    }

    /** Public lookup by name, repeated; each lookup builds a fresh parameter array. */
    static int lookupKernel(int ops, int seed) throws Exception {
        int acc = seed;
        int i = 0;
        while (i < ops) {
            Method method = ReflectTarget.class.getMethod("combine",
                    new Class[]{Integer.TYPE, Integer.TYPE});
            acc = acc ^ (method == null ? 17 : 31);
            i = i + 1;
        }
        return acc;
    }

    /**
     * Repeated invocation of one cached Method. The argument array is reused;
     * boxing is part of the operation, as with any reflective call.
     */
    static int invokeKernel(Method method, Object[] args, int ops, int seed) throws Exception {
        int acc = seed;
        int i = 0;
        while (i < ops) {
            args[0] = Integer.valueOf(acc & 1023);
            args[1] = Integer.valueOf(i & 255);
            Object result = method.invoke(null, args);
            acc = acc + ((Integer) result).intValue();
            i = i + 1;
        }
        return acc;
    }
}
