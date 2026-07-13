// Category: class initialization — recursive same-context request.
// Reading x while RecursiveInit is already initializing must complete normally
// without running <clinit> a second time. The in-progress field still has its
// default value. Expected: "init", "0", "5".
class RecursiveInit {
    static int x = initialize();
    static int initialize() {
        System.out.println("init");
        System.out.println(RecursiveInit.x);
        return 5;
    }
}
public class RecursiveInitialization {
    public static void main(String[] args) {
        System.out.println(RecursiveInit.x);
    }
}
