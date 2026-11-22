public final class Faults {
    static final class Signal extends RuntimeException {}

    static final class Bomb {
        int value;

        void arm() {
            value = 1;
        }

        void fire() {
            throw new Signal();
        }
    }

    static void covered(Bomb bomb) {
        try {
            bomb.fire();
        } catch (Signal e) {
            System.out.println(14);
        } catch (RuntimeException e) {
            System.out.println(-14);
        }
    }

    // The harness moves end_pc to the fire invocation for the boundary case.
    static void boundary(Bomb bomb) {
        try {
            bomb.arm();
            bomb.fire();
        } catch (Signal e) {
            System.out.println(15);
        }
    }

    public static void main(String[] args) {
        Bomb bomb = null;
        try {
            System.out.println(bomb.value);
        } catch (NullPointerException e) {
            System.out.println(11);
        }

        Object[] array = new Bomb[1];
        try {
            Object element = array[1];
            array[0] = element;
        } catch (ArrayIndexOutOfBoundsException e) {
            System.out.println(12);
        }

        try {
            array[0] = new Object();
        } catch (ArrayStoreException e) {
            System.out.println(13);
        }

        bomb = new Bomb();
        covered(bomb);
        try {
            boundary(bomb);
        } catch (Signal e) {
            System.out.println(16);
        }
        System.out.println(bomb.value);
    }
}
