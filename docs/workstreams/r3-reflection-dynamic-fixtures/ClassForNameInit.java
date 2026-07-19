public class ClassForNameInit {
    static class Target {
        static int value = initialize();

        static int initialize() {
            System.out.println("initialized");
            return 42;
        }
    }

    public static void main(String[] args) throws Exception {
        System.out.println("before");
        Class<?> first = Class.forName("ClassForNameInit$Target", false,
                ClassForNameInit.class.getClassLoader());
        System.out.println(first.getName());
        System.out.println("trigger");
        Class<?> second = Class.forName("ClassForNameInit$Target");
        System.out.println(first == second);
        System.out.println(Target.value);
    }
}
