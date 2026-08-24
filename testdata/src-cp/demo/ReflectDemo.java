// Reflection demo fixture (P-0011): Class.forName, declared members,
// Constructor.newInstance, Field.get/set, Method.invoke — public surface
// only, so output is byte-comparable against the reference JVM.
package demo;

public class ReflectDemo {
    public String tag;
    public int count;

    public ReflectDemo() {}

    public String describe() {
        return "tag=" + tag + ",count=" + count;
    }

    public static void main(String[] args) throws Exception {
        Class<?> c = Class.forName("demo.ReflectDemo");
        System.out.println("class=" + c.getName());

        Object o = c.getDeclaredConstructors()[0].newInstance();
        for (java.lang.reflect.Field f : c.getDeclaredFields()) {
            System.out.println("field:" + f.getName() + ":" + f.getType().getSimpleName());
        }
        java.lang.reflect.Field tag = c.getField("tag");
        tag.set(o, "hello");
        java.lang.reflect.Field count = c.getField("count");
        count.set(o, 7);

        for (java.lang.reflect.Method m : c.getDeclaredMethods()) {
            if (m.getName().equals("describe")) {
                System.out.println((String) m.invoke(o));
            }
        }
        System.out.println("isInstance=" + c.isInstance(o));

        // primitive Class identity
        System.out.println("intClass=" + (int.class == Integer.TYPE));

        // framework-shaped usage: generic serialization of an unknown POJO
        System.out.println(demo.MiniJson.serialize(o));
    }
}
