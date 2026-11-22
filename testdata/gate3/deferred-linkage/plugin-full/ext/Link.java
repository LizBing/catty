package ext;

public class Link {
    public static int good() { return 1; }

    public static int field() { return Dep.value; }

    public static int method() { return Dep.method(); }

    public static int klass() { return Gone.make(); }
}
