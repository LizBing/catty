// Research-only probes for R2-GATE. These classes are deliberately outside
// tests/fixtures so the regular R1 e2e loop does not treat unresolved R2
// capabilities as passing or failing release fixtures.

class IntegerToStringProbe {
    public static void main(String[] args) {
        System.out.println(Integer.toString(123456789));
        System.out.println(Integer.toString(Integer.MIN_VALUE));
    }
}

class LongToStringProbe {
    public static void main(String[] args) {
        System.out.println(Long.toString(1234567890123456789L));
        System.out.println(Long.toString(Long.MIN_VALUE));
    }
}

class DoubleParseProbe {
    public static void main(String[] args) {
        double a = Double.parseDouble("3.5");
        double b = Double.parseDouble("-0.0");
        System.out.println(a == 3.5);
        System.out.println(Double.doubleToRawLongBits(b));
    }
}

class HashMapProbe {
    public static void main(String[] args) {
        java.util.HashMap<String, Integer> map = new java.util.HashMap<>();
        map.put("alpha", 1);
        map.put("beta", 2);
        map.put(null, 3);
        System.out.println(map.size());
        System.out.println(map.get("alpha"));
        System.out.println(map.get(null));
        System.out.println(map.remove("beta"));
    }
}
