// R-0008 spike A: emitted Java code calling embedder-registered Go
// functions through the synthesized-class path.
public class GoCall {
    public static void main(String[] args) {
        System.out.println(go.Bridge.add(20, 22));
        System.out.println(go.Bridge.greet("catty"));
        try {
            go.Bridge.fail();
            System.out.println("unreachable");
        } catch (RuntimeException e) {
            System.out.println("caught:" + e.getMessage());
        }
    }
}
