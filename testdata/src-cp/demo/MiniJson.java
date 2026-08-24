// Mini-serializer: the framework-shaped proof that the reflection minimal
// surface enables real patterns — enumerate declared fields of an unknown
// POJO via Class/Field, read values generically, emit JSON-ish output.
// No direct references to Pojo types anywhere in this file.
package demo;

import java.lang.reflect.Field;

public class MiniJson {

    public static String serialize(Object o) throws Exception {
        Class<?> c = o.getClass();
        StringBuilder sb = new StringBuilder("{");
        boolean first = true;
        for (Field f : c.getDeclaredFields()) {
            if (!first) sb.append(',');
            first = false;
            sb.append('"').append(f.getName()).append("\":");
            Object v = f.get(o);
            if (v instanceof String) {
                sb.append('"').append(v).append('"');
            } else {
                sb.append(v);
            }
        }
        return sb.append('}').toString();
    }
}
