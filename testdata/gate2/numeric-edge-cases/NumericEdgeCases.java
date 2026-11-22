public final class NumericEdgeCases {
    private static final long MASK_HI = 0xFF00FF00FF00FF00L;
    private static final long MASK_LO = 0x00FF00FF00FF00FFL;

    int intField;
    long longField;
    int[] intArray;
    long[] longArray;
    Object[] objectArray;
    Object objectField;

    static final class Marker extends RuntimeException {
        private static final long serialVersionUID = 1L;

        Marker() {
            super();
        }

        Marker(Throwable cause) {
            super(cause);
        }
    }

    static int leaf(int value) {
        return value * 3 + 1;
    }

    static int intShiftCheck(int value, int distance) {
        int shl0 = value << 0;
        int shl1 = value << 1;
        int shl31 = value << 31;
        int shl32 = value << 32;
        int shl33 = value << 33;
        int shlNeg = value << -1;
        int shr0 = value >> 0;
        int shr1 = value >> 1;
        int shr31 = value >> 31;
        int shr32 = value >> 32;
        int shr33 = value >> 33;
        int shrNeg = value >> -1;
        int ushr0 = value >>> 0;
        int ushr1 = value >>> 1;
        int ushr31 = value >>> 31;
        int ushr32 = value >>> 32;
        int ushr33 = value >>> 33;
        int ushrNeg = value >>> -1;
        int mixed = shl0 ^ shl1 ^ shl31 ^ shl32 ^ shl33 ^ shlNeg
                ^ shr0 ^ shr1 ^ shr31 ^ shr32 ^ shr33 ^ shrNeg
                ^ ushr0 ^ ushr1 ^ ushr31 ^ ushr32 ^ ushr33 ^ ushrNeg;
        mixed += distance;
        System.out.println(mixed);
        return mixed;
    }

    static int longShiftCheck(long value, int distance) {
        long shl0 = value << 0;
        long shl1 = value << 1;
        long shl63 = value << 63;
        long shl64 = value << 64;
        long shl65 = value << 65;
        long shr0 = value >> 0;
        long shr1 = value >> 1;
        long shr63 = value >> 63;
        long shr64 = value >> 64;
        long shr65 = value >> 65;
        long ushr0 = value >>> 0;
        long ushr1 = value >>> 1;
        long ushr63 = value >>> 63;
        long ushr64 = value >>> 64;
        long ushr65 = value >>> 65;
        long mixed = shl0 ^ shl1 ^ shl63 ^ shl64 ^ shl65
                ^ shr0 ^ shr1 ^ shr63 ^ shr64 ^ shr65
                ^ ushr0 ^ ushr1 ^ ushr63 ^ ushr64 ^ ushr65;
        int acc = (int) mixed + (int) (mixed >>> 32) + distance;
        acc += intShiftCheck((int) value, distance);
        System.out.println(acc);
        return acc;
    }

    static int longOps(long seed, int flag) {
        long mask = MASK_HI ^ MASK_LO;
        long a = seed ^ mask;
        long b = 0x0102030405060708L;
        long sum = a + b;
        long diff = a - b;
        long product = a * b;
        long divisor = (b | 1L) + (flag & 1);
        long quotient = sum / divisor;
        long negated = -diff;
        long anded = sum & MASK_HI;
        long ored = diff | MASK_LO;
        long xored = product ^ negated;
        long minValue = Long.MIN_VALUE;
        long negativeOne = -1L;
        long minDiv = minValue / negativeOne;
        long maxValue = Long.MAX_VALUE;
        long one = 1L;
        long overflow = maxValue + one;
        long underflow = minValue - one;
        int acc = (int) sum;
        acc ^= (int) (sum >>> 32);
        acc ^= (int) diff;
        acc ^= (int) (diff >>> 32);
        acc ^= (int) product;
        acc ^= (int) (product >>> 32);
        acc ^= (int) quotient;
        acc ^= (int) (quotient >>> 32);
        acc ^= (int) anded;
        acc ^= (int) ored;
        acc ^= (int) xored;
        acc ^= (int) minDiv;
        acc ^= (int) (minDiv >>> 32);
        acc ^= (int) overflow;
        acc ^= (int) underflow;
        System.out.println(sum);
        System.out.println(product);
        System.out.println(minDiv);
        System.out.println(overflow);
        acc += longShiftCheck(a, flag);
        System.out.println(acc);
        return acc;
    }

    static int intArith(int a, int b) {
        int add = a + b;
        int sub = a - b;
        int mul = a * b;
        int anded = a & b;
        int ored = a | b;
        int xored = a ^ b;
        int negated = -a;
        int minValue = Integer.MIN_VALUE;
        int negativeOne = -1;
        int minDiv = minValue / negativeOne;
        int zero = 0;
        int quotient = 0;
        try {
            quotient = a / zero;
            System.out.println(-7001);
        } catch (ArithmeticException error) {
            System.out.println(7001);
        }
        int acc = add ^ sub ^ mul ^ anded ^ ored ^ xored ^ negated ^ minDiv ^ quotient;
        acc += leaf(a) + leaf(b);
        System.out.println(acc);
        return acc;
    }

    static int iincDemo(int start) {
        int counter = start;
        counter += 7;
        counter -= 3;
        counter += 100;
        counter -= 250;
        counter += 1;
        counter -= 1;
        System.out.println(counter);
        return counter;
    }

    static int constantKinds() {
        int c0 = 0;
        int c1 = 1;
        int c2 = 2;
        int c3 = 3;
        int c4 = 4;
        int c5 = 5;
        int cm1 = -1;
        int bytePush = 100;
        int shortPush = 30000;
        int intPush = 100000;
        int sum = c0 + c1 + c2 + c3 + c4 + c5 + cm1 + bytePush + shortPush + intPush;
        long high = 0x1122334455667788L;
        long low = 0x8877665544332211L;
        long mixed = high ^ low;
        long pool = 0;
        pool += 100000;
        pool += 100001;
        pool += 100002;
        pool += 100003;
        pool += 100004;
        pool += 100005;
        pool += 100006;
        pool += 100007;
        pool += 100008;
        pool += 100009;
        pool += 100010;
        pool += 100011;
        pool += 100012;
        pool += 100013;
        pool += 100014;
        pool += 100015;
        pool += 100016;
        pool += 100017;
        pool += 100018;
        pool += 100019;
        pool += 100020;
        pool += 100021;
        pool += 100022;
        pool += 100023;
        pool += 100024;
        pool += 100025;
        pool += 100026;
        pool += 100027;
        pool += 100028;
        pool += 100029;
        pool += 100030;
        pool += 100031;
        pool += 100032;
        pool += 100033;
        pool += 100034;
        pool += 100035;
        pool += 100036;
        pool += 100037;
        pool += 100038;
        pool += 100039;
        pool += 100040;
        pool += 100041;
        pool += 100042;
        pool += 100043;
        pool += 100044;
        pool += 100045;
        pool += 100046;
        pool += 100047;
        pool += 100048;
        pool += 100049;
        pool += 100050;
        pool += 100051;
        pool += 100052;
        pool += 100053;
        pool += 100054;
        pool += 100055;
        pool += 100056;
        pool += 100057;
        pool += 100058;
        pool += 100059;
        pool += 100060;
        pool += 100061;
        pool += 100062;
        pool += 100063;
        pool += 100064;
        pool += 100065;
        pool += 100066;
        pool += 100067;
        pool += 100068;
        pool += 100069;
        pool += 100070;
        pool += 100071;
        pool += 100072;
        pool += 100073;
        pool += 100074;
        pool += 100075;
        pool += 100076;
        pool += 100077;
        pool += 100078;
        pool += 100079;
        pool += 100080;
        pool += 100081;
        pool += 100082;
        pool += 100083;
        pool += 100084;
        pool += 100085;
        pool += 100086;
        pool += 100087;
        pool += 100088;
        pool += 100089;
        pool += 100090;
        pool += 100091;
        pool += 100092;
        pool += 100093;
        pool += 100094;
        pool += 100095;
        pool += 100096;
        pool += 100097;
        pool += 100098;
        pool += 100099;
        pool += 100100;
        pool += 100101;
        pool += 100102;
        pool += 100103;
        pool += 100104;
        pool += 100105;
        pool += 100106;
        pool += 100107;
        pool += 100108;
        pool += 100109;
        pool += 100110;
        pool += 100111;
        pool += 100112;
        pool += 100113;
        pool += 100114;
        pool += 100115;
        pool += 100116;
        pool += 100117;
        pool += 100118;
        pool += 100119;
        pool += 100120;
        pool += 100121;
        pool += 100122;
        pool += 100123;
        pool += 100124;
        pool += 100125;
        pool += 100126;
        pool += 100127;
        pool += 100128;
        pool += 100129;
        pool += 100130;
        pool += 100131;
        pool += 100132;
        pool += 100133;
        pool += 100134;
        pool += 100135;
        pool += 100136;
        pool += 100137;
        pool += 100138;
        pool += 100139;
        pool += 100140;
        pool += 100141;
        pool += 100142;
        pool += 100143;
        pool += 100144;
        pool += 100145;
        pool += 100146;
        pool += 100147;
        pool += 100148;
        pool += 100149;
        pool += 100150;
        pool += 100151;
        pool += 100152;
        pool += 100153;
        pool += 100154;
        pool += 100155;
        pool += 100156;
        pool += 100157;
        pool += 100158;
        pool += 100159;
        pool += 100160;
        pool += 100161;
        pool += 100162;
        pool += 100163;
        pool += 100164;
        pool += 100165;
        pool += 100166;
        pool += 100167;
        pool += 100168;
        pool += 100169;
        pool += 100170;
        pool += 100171;
        pool += 100172;
        pool += 100173;
        pool += 100174;
        pool += 100175;
        pool += 100176;
        pool += 100177;
        pool += 100178;
        pool += 100179;
        pool += 100180;
        pool += 100181;
        pool += 100182;
        pool += 100183;
        pool += 100184;
        pool += 100185;
        pool += 100186;
        pool += 100187;
        pool += 100188;
        pool += 100189;
        pool += 100190;
        pool += 100191;
        pool += 100192;
        pool += 100193;
        pool += 100194;
        pool += 100195;
        pool += 100196;
        pool += 100197;
        pool += 100198;
        pool += 100199;
        pool += 100200;
        pool += 100201;
        pool += 100202;
        pool += 100203;
        pool += 100204;
        pool += 100205;
        pool += 100206;
        pool += 100207;
        pool += 100208;
        pool += 100209;
        pool += 100210;
        pool += 100211;
        pool += 100212;
        pool += 100213;
        pool += 100214;
        pool += 100215;
        pool += 100216;
        pool += 100217;
        pool += 100218;
        pool += 100219;
        pool += 100220;
        pool += 100221;
        pool += 100222;
        pool += 100223;
        pool += 100224;
        pool += 100225;
        pool += 100226;
        pool += 100227;
        pool += 100228;
        pool += 100229;
        pool += 100230;
        pool += 100231;
        pool += 100232;
        pool += 100233;
        pool += 100234;
        pool += 100235;
        pool += 100236;
        pool += 100237;
        pool += 100238;
        pool += 100239;
        pool += 100240;
        pool += 100241;
        pool += 100242;
        pool += 100243;
        pool += 100244;
        pool += 100245;
        pool += 100246;
        pool += 100247;
        pool += 100248;
        pool += 100249;
        pool += 100250;
        pool += 100251;
        pool += 100252;
        pool += 100253;
        pool += 100254;
        pool += 100255;
        pool += 100256;
        pool += 100257;
        pool += 100258;
        pool += 100259;
        pool += 100260;
        pool += 100261;
        pool += 100262;
        pool += 100263;
        pool += 100264;
        pool += 100265;
        pool += 100266;
        pool += 100267;
        pool += 100268;
        pool += 100269;
        pool += 100270;
        pool += 100271;
        pool += 100272;
        pool += 100273;
        pool += 100274;
        pool += 100275;
        pool += 100276;
        pool += 100277;
        pool += 100278;
        pool += 100279;
        pool += 100280;
        pool += 100281;
        pool += 100282;
        pool += 100283;
        pool += 100284;
        pool += 100285;
        pool += 100286;
        pool += 100287;
        pool += 100288;
        pool += 100289;
        pool += 100290;
        pool += 100291;
        pool += 100292;
        pool += 100293;
        pool += 100294;
        pool += 100295;
        pool += 100296;
        pool += 100297;
        pool += 100298;
        pool += 100299;
        int acc = sum;
        acc += (int) mixed;
        acc += (int) (mixed >>> 32);
        acc += (int) pool;
        acc += (int) (pool >>> 32);
        System.out.println(acc);
        return acc;
    }

    static int convWiden(int i, long l) {
        long i2l = i;
        float i2f = i;
        double i2d = i;
        int l2i = (int) l;
        float l2f = l;
        double l2d = l;
        int acc = (int) i2l;
        acc ^= (int) (i2l >>> 32);
        acc ^= l2i;
        acc ^= Float.floatToRawIntBits(i2f);
        long i2dBits = Double.doubleToRawLongBits(i2d);
        acc ^= (int) i2dBits;
        acc ^= (int) (i2dBits >>> 32);
        acc ^= Float.floatToRawIntBits(l2f);
        long l2dBits = Double.doubleToRawLongBits(l2d);
        acc ^= (int) l2dBits;
        acc ^= (int) (l2dBits >>> 32);
        System.out.println(acc);
        return acc;
    }

    static int convFloat(float f, double d) {
        double f2d = f;
        float d2f = (float) d;
        long fBits = Double.doubleToRawLongBits(f2d);
        int dBits = Float.floatToRawIntBits(d2f);
        System.out.println(fBits);
        System.out.println(dBits);
        return (int) fBits ^ (int) (fBits >>> 32) ^ dBits;
    }

    static int saturation(float f, double d) {
        int f2i = (int) f;
        long f2l = (long) f;
        int d2i = (int) d;
        long d2l = (long) d;
        System.out.println(f2i);
        System.out.println(f2l);
        System.out.println(d2i);
        System.out.println(d2l);
        return f2i ^ (int) f2l ^ d2i ^ (int) d2l ^ (int) (d2l >>> 32);
    }

    static int narrowing(int value) {
        byte asByte = (byte) value;
        char asChar = (char) value;
        short asShort = (short) value;
        int acc = asByte * 31 + asChar * 17 + asShort * 7;
        System.out.println(acc);
        return acc;
    }

    static float specialFloat(int selector) {
        float zero = 0.0f;
        float one = 1.0f;
        if (selector == 0) {
            return zero / zero;
        }
        if (selector == 1) {
            return one / zero;
        }
        return -one / zero;
    }

    static double specialDouble(int selector) {
        double zero = 0.0;
        double one = 1.0;
        if (selector == 0) {
            return zero / zero;
        }
        if (selector == 1) {
            return one / zero;
        }
        return -one / zero;
    }

    static int floatOps(float a, float b) {
        float add = a + b;
        float sub = a - b;
        float mul = a * b;
        float div = a / b;
        float neg = -a;
        System.out.println(Float.floatToRawIntBits(add));
        System.out.println(Float.floatToRawIntBits(sub));
        System.out.println(Float.floatToRawIntBits(mul));
        System.out.println(Float.floatToRawIntBits(div));
        System.out.println(Float.floatToRawIntBits(neg));
        int cmp = compareFloat(a, b);
        System.out.println(cmp);
        return Float.floatToRawIntBits(add) ^ cmp;
    }

    static int compareFloat(float a, float b) {
        int bits = 0;
        if (a < b) {
            bits += 1;
        }
        if (a <= b) {
            bits += 2;
        }
        if (a > b) {
            bits += 4;
        }
        if (a >= b) {
            bits += 8;
        }
        if (a == b) {
            bits += 16;
        }
        if (a != b) {
            bits += 32;
        }
        if (a == a) {
            bits += 64;
        }
        if (b == b) {
            bits += 128;
        }
        System.out.println(bits);
        return bits;
    }

    static int doubleOps(double a, double b) {
        double add = a + b;
        double sub = a - b;
        double mul = a * b;
        double div = a / b;
        double neg = -a;
        System.out.println(Double.doubleToRawLongBits(add));
        System.out.println(Double.doubleToRawLongBits(sub));
        System.out.println(Double.doubleToRawLongBits(mul));
        System.out.println(Double.doubleToRawLongBits(div));
        System.out.println(Double.doubleToRawLongBits(neg));
        int cmp = compareDouble(a, b);
        System.out.println(cmp);
        return (int) Double.doubleToRawLongBits(add) ^ cmp;
    }

    static int compareDouble(double a, double b) {
        int bits = 0;
        if (a < b) {
            bits += 1;
        }
        if (a <= b) {
            bits += 2;
        }
        if (a > b) {
            bits += 4;
        }
        if (a >= b) {
            bits += 8;
        }
        if (a == b) {
            bits += 16;
        }
        if (a != b) {
            bits += 32;
        }
        if (a == a) {
            bits += 64;
        }
        if (b == b) {
            bits += 128;
        }
        System.out.println(bits);
        return bits;
    }

    static int nanPaths(float fn, double dn) {
        float fsum = fn + fn;
        float fsub = fn - fn;
        float fmul = fn * fn;
        float fdiv = fn / fn;
        double dsum = dn + dn;
        double dsub = dn - dn;
        double dmul = dn * dn;
        double ddiv = dn / dn;
        int acc = 0;
        if (fsum != fsum) {
            acc += 1;
        }
        if (fsub != fsub) {
            acc += 2;
        }
        if (fmul != fmul) {
            acc += 4;
        }
        if (fdiv != fdiv) {
            acc += 8;
        }
        if (dsum != dsum) {
            acc += 16;
        }
        if (dsub != dsub) {
            acc += 32;
        }
        if (dmul != dmul) {
            acc += 64;
        }
        if (ddiv != ddiv) {
            acc += 128;
        }
        System.out.println(acc);
        return acc;
    }

    static int bitRoundTrip(int floatBits, long doubleBits) {
        float f = Float.intBitsToFloat(floatBits);
        double d = Double.longBitsToDouble(doubleBits);
        int backFloat = Float.floatToRawIntBits(f);
        long backDouble = Double.doubleToRawLongBits(d);
        System.out.println(backFloat);
        System.out.println(backDouble);
        return backFloat ^ (int) backDouble ^ (int) (backDouble >>> 32);
    }

    static int longCompare(long a, long b) {
        int bits = 0;
        if (a < b) {
            bits += 1;
        }
        if (a == b) {
            bits += 2;
        }
        if (a > b) {
            bits += 4;
        }
        System.out.println(bits);
        return bits;
    }

    static int branchKinds(int a, int b) {
        int bits = 0;
        if (a < b) {
            bits += 1;
        }
        if (a >= b) {
            bits += 2;
        }
        if (a > b) {
            bits += 4;
        }
        if (a <= b) {
            bits += 8;
        }
        if (a == b) {
            bits += 16;
        }
        if (a != b) {
            bits += 32;
        }
        if (a != 0) {
            bits += 64;
        }
        if (a == 0) {
            bits += 128;
        }
        if (a < 0) {
            bits += 256;
        }
        if (a >= 0) {
            bits += 512;
        }
        if (a > 0) {
            bits += 1024;
        }
        if (a <= 0) {
            bits += 2048;
        }
        System.out.println(bits);
        return bits;
    }

    static int refBranches(Object first, Object second) {
        int bits = 0;
        if (first == second) {
            bits += 1;
        }
        if (first != second) {
            bits += 2;
        }
        if (first == null) {
            bits += 4;
        }
        if (first != null) {
            bits += 8;
        }
        if (second == null) {
            bits += 16;
        }
        if (second != null) {
            bits += 32;
        }
        System.out.println(bits);
        return bits;
    }

    static int leafUse(int value) {
        int first = leaf(value);
        int second = leaf(value + 17);
        int acc = first * 31 + second;
        System.out.println(acc);
        return acc;
    }

    static int exceptionPath(int mode) {
        Marker saved = new Marker();
        int acc = 0;
        try {
            if (mode == 0) {
                throw new Marker(saved);
            }
            acc = 5;
        } catch (Marker error) {
            int identity = (error.getCause() == saved) ? 1 : 0;
            int self = (error.getCause() == error) ? 1 : 0;
            System.out.println(1000 + identity * 10 + self);
            acc = 3;
        }
        Marker second = new Marker();
        try {
            throw second;
        } catch (Marker error) {
            System.out.println(error == second ? 2001 : -2001);
            System.out.println(error.getCause() == null ? 2002 : -2002);
        }
        System.out.println(acc);
        return acc;
    }

    static int extraPath(int code) {
        long factor = 0x100000001L;
        long mixed = factor * code;
        int acc = (int) mixed ^ (int) (mixed >>> 32) ^ 0x55aa55aa;
        System.out.println(acc);
        return acc;
    }

    static int remainders(long a, int b, float f, double d) {
        long l1 = a % 3L;
        long l2 = a % -7L;
        float f1 = f % 2.5f;
        float f2 = f % -0.0f;
        double d1 = d % 3.25;
        double d2 = d % 0.0;
        System.out.println((int) l1);
        System.out.println((int) l2);
        System.out.println(Float.floatToRawIntBits(f1));
        System.out.println((f2 == f2) ? 0 : 1);
        System.out.println(Double.doubleToRawLongBits(d1));
        System.out.println((d2 == d2) ? 0 : 1);
        int zero = 0;
        try {
            long bad = a / (long) zero;
            System.out.println((int) bad);
        } catch (ArithmeticException expected) {
            System.out.println(31337);
        }
        try {
            long bad = a % (long) zero;
            System.out.println((int) bad);
        } catch (ArithmeticException expected) {
            System.out.println(31338);
        }
        return (int) (l1 + l2);
    }

    long longMark() {
        return 0x1122334455667788L;
    }

    int stackForms(int i, int v, long lv, Object ov) {
        objectArray = new Object[2];
        intArray = new int[2];
        longArray = new long[2];
        objectField = objectArray[i] = ov;
        intField = intArray[i] = v;
        longField = longArray[i] = lv;
        objectArray[i] = objectField = ov;
        intArray[i] = intField = v + 1;
        longArray[i] = longField = lv + 2L;
        intArray[i] <<= 3;
        longArray[i] = longArray[i] + lv;
        intField = intArray[i]++;
        longField = longArray[i]++;
        System.out.println(intField);
        System.out.println((int) longField);
        System.out.println(intArray[i]);
        System.out.println((int) longArray[i]);
        System.out.println(objectArray[i] == ov ? 1 : 0);
        longArray[i] = longArray[i] + lv;
        System.out.println((int) longArray[i]);
        longMark();
        return intField;
    }

    static int unusedPath(int value) {
        return value * 13 + 99;
    }

    public static void main(String[] args) {
        int extra = args.length;
        int a = longOps(0x1122334455667788L, extra);
        System.out.println(a);
        int b = intArith(0x71234567, -37);
        System.out.println(b);
        int c = iincDemo(extra + 10);
        System.out.println(c);
        int d = constantKinds();
        System.out.println(d);
        int e = convWiden(0x12345678, 0x1122334455667788L);
        System.out.println(e);
        int f = convFloat(1.5f, -2.25);
        System.out.println(f);
        saturation(specialFloat(0), specialDouble(0));
        saturation(specialFloat(1), specialDouble(1));
        saturation(specialFloat(2), specialDouble(2));
        saturation(3.0e9f, -2147483649.0);
        saturation(2147483647f, 9223372036854775807d);
        saturation(-2147483648f, -9223372036854775808d);
        floatOps(1.5f, -2.5f);
        floatOps(-0.0f, 1.0f);
        floatOps(1.0f, -0.0f);
        floatOps(-1.0f, 0.0f);
        System.out.println(compareFloat(specialFloat(0), specialFloat(0)));
        System.out.println(compareFloat(specialFloat(0), 1.0f));
        System.out.println(compareFloat(1.0f, specialFloat(0)));
        System.out.println(compareFloat(0.0f, -0.0f));
        System.out.println(compareFloat(specialFloat(1), specialFloat(2)));
        System.out.println(compareFloat(specialFloat(1), specialFloat(1)));
        doubleOps(1.5, -2.5);
        doubleOps(-0.0, 1.0);
        doubleOps(1.0, -0.0);
        doubleOps(1.0e308, 1.0e-308);
        System.out.println(compareDouble(specialDouble(0), specialDouble(0)));
        System.out.println(compareDouble(specialDouble(0), 1.0));
        System.out.println(compareDouble(1.0, specialDouble(0)));
        System.out.println(compareDouble(0.0, -0.0));
        System.out.println(compareDouble(specialDouble(1), specialDouble(2)));
        System.out.println(compareDouble(specialDouble(1), specialDouble(1)));
        nanPaths(specialFloat(0), specialDouble(0));
        bitRoundTrip(0x7fc00001, 0x7ff8000000000001L);
        bitRoundTrip(-2147483648, Long.MIN_VALUE);
        longCompare(0x1122334455667788L, 0x1122334455667788L);
        longCompare(-5L, 9L);
        longCompare(9L, -5L);
        branchKinds(7, 3);
        branchKinds(-2, 0);
        refBranches(args, null);
        refBranches(args, args);
        refBranches(null, null);
        narrowing(0x12345678);
        leafUse(extra + 2);
        System.out.println(remainders(0x1122334455667788L, -37, 1.5f, -2.25));
        NumericEdgeCases forms = new NumericEdgeCases();
        System.out.println(forms.stackForms(1, 7, 0x1122334455667788L, args));
        exceptionPath(0);
        if (extra == 3) {
            extraPath(extra);
        }
    }
}
