// Business rules for the embeddemo: Java logic consuming host-provided
// Go capabilities (crypto + HTTP via the local server the host started).
package demo;

public class Rules {
    public static void run() {
        String digest = go.Demo.md5Hex("catty");
        System.out.println("md5(catty)=" + digest);

        long n = go.Demo.fetchLen(go.Demo.baseUrl() + "/payload");
        System.out.println("payload_bytes=" + n);

        // Rules can mix Go capabilities with plain Java logic.
        String line = "md5:" + go.Demo.md5Hex(digest);
        System.out.println("line_len=" + line.length());
    }
}
