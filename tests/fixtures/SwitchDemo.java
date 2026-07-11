public class SwitchDemo {
    static int denseSwitch(int n) {
        switch (n) {            // compiles to tableswitch (dense keys)
            case 0: return 100;
            case 1: return 101;
            case 2: return 102;
            case 3: return 103;
            case 4: return 104;
            default: return -1;
        }
    }

    static int sparseSwitch(int n) {
        switch (n) {            // compiles to lookupswitch (sparse keys)
            case 1: return 1;
            case 10: return 10;
            case 100: return 100;
            case 1000: return 1000;
            default: return 0;
        }
    }

    public static void main(String[] args) {
        for (int i = 0; i <= 5; i++) {
            System.out.println(denseSwitch(i));
        }
        System.out.println(sparseSwitch(10));
        System.out.println(sparseSwitch(1000));
        System.out.println(sparseSwitch(7));
    }
}
