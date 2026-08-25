// Reflection v2 fixture (P-0012): inherited member traversal.
package demo;

class IBase {
    public String base = "B";
    private int secret = 0;

    public String baseInfo() { return "base:" + base; }
}

public class InheritDemo extends IBase {
    public int num = 7;
    public String base = "C"; // shadows IBase.base (distinct fields)

    public String baseInfo() { return "child:" + base; }

    public static void main(String[] args) throws Exception {
        Class<?> c = Class.forName("demo.InheritDemo");

        java.lang.reflect.Field[] fs = c.getFields();
        System.out.println("count=" + fs.length);
        for (java.lang.reflect.Field f : fs) {
            System.out.println("field=" + f.getDeclaringClass().getSimpleName()
                + "." + f.getName());
        }

        InheritDemo d = new InheritDemo();
        d.base = "X";
        Object o = d;
        java.lang.reflect.Method m = c.getMethod("baseInfo");
        System.out.println((String) m.invoke(o));

        java.lang.reflect.Method inherited = o.getClass().getMethod("baseInfo");
        System.out.println("same=" + m.equals(inherited));

        Class<?> sup = c.getSuperclass();
        System.out.println("super=" + sup.getSimpleName());
        System.out.println("done");
    }
}
