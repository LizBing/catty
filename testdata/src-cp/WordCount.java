import java.util.HashMap;
import java.util.HashSet;

public class WordCount {
    public static void main(String[] args) {
        String text = "the quick brown fox jumps over the lazy dog the end";
        String[] words = text.split(" ");

        HashMap counts = new HashMap();
        HashSet unique = new HashSet();

        for (int i = 0; i < words.length; i++) {
            String w = words[i];
            Integer c = (Integer) counts.get(w);
            int n;
            if (c == null) {
                n = 0;
            } else {
                n = c.intValue();
            }
            counts.put(w, Integer.valueOf(n + 1));
            unique.add(w);
        }

        System.out.println("total words: " + words.length);
        System.out.println("unique words: " + unique.size());
        System.out.println("count(the): " + counts.get("the"));
        System.out.println("count(fox): " + counts.get("fox"));
        System.out.println("count(cat): " + counts.get("cat"));

        // Exercise additional String methods
        String greeting = "  Hello, Catty!  ";
        String trimmed = greeting.trim();
        System.out.println("trimmed: [" + trimmed + "]");
        System.out.println("upper: " + trimmed.toUpperCase());
        System.out.println("lower: " + trimmed.toLowerCase());
        System.out.println("substr(7,12): " + trimmed.substring(7, 12));
        System.out.println("contains(Catty): " + trimmed.contains("Catty"));
        System.out.println("startsWith(Hel): " + trimmed.startsWith("Hel"));
        System.out.println("isEmpty: " + "".isEmpty());

        // Long wrapper
        Long big = Long.valueOf(9999999999L);
        long val = big.longValue();
        System.out.println("long value: " + val);

        Boolean flag = Boolean.valueOf(true);
        System.out.println("boolean: " + flag.booleanValue() + " parsed=" + Boolean.parseBoolean("true"));

        System.out.println("done");
    }
}
