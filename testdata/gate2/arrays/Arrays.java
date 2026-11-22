public final class Arrays {
    interface Sized {
        int size();
    }

    static final class Item {
        final int id;

        Item(int id) {
            this.id = id;
        }
    }

    static final class IntSized implements Sized {
        private final int[] values;

        IntSized(int[] values) {
            this.values = values;
        }

        public int size() {
            return values.length;
        }
    }

    static final class Denied extends RuntimeException {
        private static final long serialVersionUID = 1L;
    }

    static int lengthOf(int[] values) {
        return values.length;
    }

    static int deny() {
        throw new Denied();
    }

    static void mutate(int[] values) {
        values[0] = 55;
    }

    static int[] replace(int[] values) {
        values = new int[] {7};
        return values;
    }

    static void defaults() {
        int[] ints = new int[3];
        long[] longs = new long[2];
        boolean[] booleans = new boolean[2];
        char[] chars = new char[2];
        byte[] bytes = new byte[2];
        short[] shorts = new short[2];
        float[] floats = new float[2];
        double[] doubles = new double[2];
        System.out.println(10);
        System.out.println(ints.length);
        System.out.println(ints[0]);
        System.out.println(longs.length);
        System.out.println((int) longs[0]);
        System.out.println(booleans.length);
        System.out.println(booleans[0] ? 1 : 0);
        System.out.println(chars.length);
        System.out.println((int) chars[0]);
        System.out.println(bytes.length);
        System.out.println((int) bytes[0]);
        System.out.println(shorts.length);
        System.out.println((int) shorts[0]);
        System.out.println(floats.length);
        System.out.println((int) floats[0]);
        System.out.println(doubles.length);
        System.out.println((int) doubles[0]);
        System.out.println(lengthOf(ints));
        System.out.println(lengthOf(new int[0]));
    }

    static void primitives() {
        int[] ints = new int[3];
        ints[0] = 11;
        ints[1] = 22;
        ints[2] = 33;
        long[] longs = new long[2];
        longs[0] = 0x100000000L;
        longs[1] = 5L;
        boolean[] booleans = new boolean[2];
        booleans[1] = true;
        char[] chars = new char[2];
        chars[0] = 'A';
        chars[1] = 'Z';
        byte[] bytes = new byte[2];
        bytes[0] = (byte) 7;
        bytes[1] = (byte) -3;
        short[] shorts = new short[2];
        shorts[0] = (short) 1000;
        shorts[1] = (short) -1000;
        float[] floats = new float[2];
        floats[0] = 1.5f;
        floats[1] = 2.5f;
        double[] doubles = new double[2];
        doubles[0] = 1.75;
        doubles[1] = 2.25;
        System.out.println(20);
        System.out.println(ints[2]);
        System.out.println((int) (longs[0] >>> 32));
        System.out.println((int) longs[1]);
        System.out.println(booleans[1] ? 1 : 0);
        System.out.println((int) chars[0]);
        System.out.println((int) chars[1]);
        System.out.println((int) bytes[0]);
        System.out.println((int) bytes[1]);
        System.out.println((int) shorts[0]);
        System.out.println((int) shorts[1]);
        System.out.println((int) floats[0]);
        System.out.println((int) floats[1]);
        System.out.println((int) doubles[0]);
        System.out.println((int) doubles[1]);
    }

    static void references() {
        Item[] items = new Item[2];
        System.out.println(30);
        System.out.println(items.length);
        System.out.println(items[0] == null ? 1 : 0);
        items[0] = new Item(41);
        items[1] = new Item(42);
        System.out.println(items[0].id);
        System.out.println(items[1].id);
        String[] strings = new String[2];
        String first = "a";
        strings[0] = first;
        Object[] objects = strings;
        String second = "b";
        objects[1] = second;
        System.out.println(strings.length);
        System.out.println(strings[1] == second ? 1 : 0);
        try {
            objects[0] = new Item(9);
            System.out.println(-31);
        } catch (ArrayStoreException expected) {
            System.out.println(31);
        }
        System.out.println(strings[0] == first ? 1 : 0);
        try {
            objects[0] = items[1];
            System.out.println(-32);
        } catch (ArrayStoreException expected) {
            System.out.println(32);
        }
        Object boxed = strings;
        System.out.println(strings instanceof Object ? 1 : 0);
        System.out.println(strings instanceof Cloneable ? 1 : 0);
        System.out.println(strings instanceof java.io.Serializable ? 1 : 0);
        Cloneable cloneable = (Cloneable) boxed;
        System.out.println(cloneable == null ? 0 : 1);
        java.io.Serializable serializable = (java.io.Serializable) boxed;
        System.out.println(serializable == null ? 0 : 1);
        String[] restored = (String[]) boxed;
        System.out.println(restored.length);
        int[] numbers = new int[2];
        Object primitive = numbers;
        System.out.println(primitive instanceof Object ? 1 : 0);
        System.out.println(primitive instanceof Cloneable ? 1 : 0);
        Cloneable primitiveCloneable = (Cloneable) primitive;
        System.out.println(primitiveCloneable == null ? 0 : 1);
        java.io.Serializable primitiveSerializable = (java.io.Serializable) primitive;
        System.out.println(primitiveSerializable == null ? 0 : 1);
        System.out.println(primitive instanceof Object[] ? 0 : 1);
        try {
            Object[] wrong = (Object[]) primitive;
            System.out.println(wrong.length);
        } catch (ClassCastException expected) {
            System.out.println(33);
        }
        Object single = new Item(50);
        Item item = (Item) single;
        System.out.println(item.id);
        Sized sized = new IntSized(numbers);
        System.out.println(sized.size());
    }

    static void multi() {
        int[][] grid = new int[2][3];
        System.out.println(40);
        System.out.println(grid.length);
        System.out.println(grid[0].length);
        System.out.println(grid[1].length);
        System.out.println(grid[0] == grid[1] ? 0 : 1);
        System.out.println(grid[0][0]);
        grid[1][2] = 9;
        System.out.println(grid[1][2]);
        int[][] partial = new int[2][];
        System.out.println(partial.length);
        System.out.println(partial[0] == null ? 1 : 0);
        System.out.println(partial[1] == null ? 1 : 0);
        partial[0] = new int[3];
        System.out.println(partial[0].length);
        System.out.println(partial[1] == null ? 1 : 0);
        int[][][] cube = new int[2][3][4];
        System.out.println(cube.length);
        System.out.println(cube[0].length);
        System.out.println(cube[0][2].length);
        System.out.println(cube[1][2][3]);
        cube[1][2][3] = 7;
        System.out.println(cube[1][2][3]);
        System.out.println(cube[0] == cube[1] ? 0 : 1);
        System.out.println(cube[0][0] == cube[0][1] ? 0 : 1);
        int[][] ragged = new int[3][];
        ragged[0] = new int[1];
        ragged[1] = new int[2];
        ragged[2] = new int[3];
        System.out.println(ragged[0].length);
        System.out.println(ragged[1].length);
        System.out.println(ragged[2].length);
        System.out.println(ragged[2][2]);
        int[] row = new int[2];
        int[][] shared = new int[2][];
        shared[0] = row;
        shared[1] = row;
        System.out.println(shared[0] == shared[1] ? 1 : 0);
        shared[0][0] = 5;
        System.out.println(shared[1][0]);
        int[][] zeroRows = new int[2][0];
        System.out.println(zeroRows.length);
        System.out.println(zeroRows[0] == null ? 0 : 1);
        System.out.println(zeroRows[0].length);
    }

    static void failures() {
        System.out.println(50);
        int[] none = null;
        try {
            System.out.println(none[0]);
        } catch (NullPointerException expected) {
            System.out.println(51);
        }
        try {
            none[0] = 5;
        } catch (NullPointerException expected) {
            System.out.println(52);
        }
        try {
            System.out.println(none.length);
        } catch (NullPointerException expected) {
            System.out.println(53);
        }
        try {
            int[] bad = new int[-1];
            System.out.println(bad.length);
        } catch (NegativeArraySizeException expected) {
            System.out.println(54);
        }
        try {
            int[][] bad = new int[2][-1];
            System.out.println(bad.length);
        } catch (NegativeArraySizeException expected) {
            System.out.println(55);
        }
        int[] ints = new int[2];
        try {
            System.out.println(ints[2]);
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(56);
        }
        try {
            ints[-1] = 5;
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(57);
        }
        long[] longs = new long[2];
        try {
            System.out.println((int) longs[3]);
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(58);
        }
        try {
            longs[2] = 1L;
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(59);
        }
        Object[] objects = new Object[2];
        System.out.println(objects[0] == null ? 1 : 0);
        try {
            System.out.println(objects[2] == null ? 0 : 1);
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(60);
        }
        try {
            objects[-1] = "x";
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(61);
        }
        Object[] covariance = new String[1];
        try {
            covariance[0] = new Item(1);
            System.out.println(-62);
        } catch (ArrayStoreException expected) {
            System.out.println(62);
        }
        System.out.println(covariance[0] == null ? 1 : 0);
        Object other = new Item(2);
        try {
            String failing = (String) other;
            System.out.println(failing.length());
        } catch (ClassCastException expected) {
            System.out.println(63);
        }
        try {
            System.out.println(none[5]);
        } catch (NullPointerException expected) {
            System.out.println(64);
        } catch (ArrayIndexOutOfBoundsException wrong) {
            System.out.println(-64);
        }
        try {
            none[5] = 7;
        } catch (NullPointerException expected) {
            System.out.println(65);
        } catch (ArrayIndexOutOfBoundsException wrong) {
            System.out.println(-65);
        }
        Object[] bounded = new String[2];
        try {
            bounded[5] = new Item(3);
        } catch (ArrayIndexOutOfBoundsException expected) {
            System.out.println(66);
        } catch (ArrayStoreException wrong) {
            System.out.println(-66);
        }
        Object primitive = new int[1];
        try {
            ((Object[]) primitive)[99] = "x";
        } catch (ClassCastException expected) {
            System.out.println(67);
        } catch (ArrayIndexOutOfBoundsException wrong) {
            System.out.println(-67);
        }
        try {
            Object[] cast = (Object[]) primitive;
            System.out.println(cast.length);
        } catch (ClassCastException expected) {
            System.out.println(68);
        }
        Object[] nullCovariance = null;
        try {
            nullCovariance[0] = new Item(6);
        } catch (NullPointerException expected) {
            System.out.println(91);
        } catch (ArrayStoreException wrong) {
            System.out.println(-91);
        }
        try {
            System.out.println(deny());
        } catch (Denied expected) {
            System.out.println(69);
        }
    }

    static void aliasing() {
        System.out.println(80);
        int[] data = new int[3];
        data[0] = 1;
        data[1] = 2;
        data[2] = 3;
        mutate(data);
        System.out.println(data[0]);
        int[] replacement = replace(data);
        System.out.println(data[0]);
        System.out.println(replacement[0]);
        System.out.println(replacement.length);
        System.out.println(lengthOf(data));
        int[][] matrix = new int[2][2];
        int[] firstRow = matrix[0];
        firstRow[1] = 8;
        System.out.println(matrix[0][1]);
        matrix[1] = firstRow;
        System.out.println(matrix[0] == matrix[1] ? 1 : 0);
        String[] words = new String[2];
        String q = "q";
        words[0] = q;
        String[] alias = words;
        String r = "r";
        alias[1] = r;
        System.out.println(words[1] == r ? 1 : 0);
    }

    static void extraCovariance() {
        System.out.println(85);
        Object[] values = new String[2];
        String left = "left";
        values[0] = left;
        System.out.println(values.length);
        System.out.println(values[0] == left ? 1 : 0);
        String right = "right";
        values[1] = right;
        System.out.println(values[1] == right ? 1 : 0);
        try {
            values[0] = new Item(4);
            System.out.println(-81);
        } catch (ArrayStoreException expected) {
            System.out.println(81);
        }
        System.out.println(values[0] == left ? 1 : 0);
        Object primitive = new long[2];
        System.out.println(primitive instanceof Object[] ? 0 : 1);
        try {
            Object[] bad = (Object[]) primitive;
            System.out.println(bad.length);
        } catch (ClassCastException expected) {
            System.out.println(82);
        }
        int[] numbers = new int[3];
        Sized sized = new IntSized(numbers);
        System.out.println(sized.size());
    }

    public static void main(String[] args) {
        defaults();
        primitives();
        references();
        multi();
        failures();
        aliasing();
        if (args.length > 0) {
            extraCovariance();
        }
    }
}
