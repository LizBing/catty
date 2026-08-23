import java.io.BufferedReader;
import java.io.IOException;
import java.io.FileReader;

public class WcFile {
    public static void main(String[] args) throws IOException {
        if (args.length == 0) {
            System.out.println("usage: wcfile <input.txt>");
            return;
        }
        java.util.HashMap counts = new java.util.HashMap();
        int lines = 0, words = 0;
        BufferedReader r = new BufferedReader(new FileReader(args[0]));
        String line = r.readLine();
        while (line != null) {
            lines++;
            String[] parts = line.trim().split("\\s+");
            for (int i = 0; i < parts.length; i++) {
                String w = parts[i];
                if (w.equals("")) continue;
                words++;
                Integer c = (Integer) counts.get(w);
                int n = (c == null) ? 0 : c.intValue();
                counts.put(w, Integer.valueOf(n + 1));
            }
            line = r.readLine();
        }
        r.close();
        System.out.println("lines: " + lines);
        System.out.println("words: " + words);
        System.out.println("unique: " + counts.size());
        System.out.println("top(the): " + counts.get("the"));
    }
}
