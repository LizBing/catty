public class BlockOnly {
    public static void main(String[] args) {
        synchronized (args) {
            System.out.println(73);
        }
    }
}
