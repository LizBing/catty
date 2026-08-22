import java.io.InputStream;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;

public class HttpEcho {
    public static void main(String[] args) throws Exception {
        int port = args.length > 0 ? Integer.parseInt(args[0]) : 18080;
        ServerSocket ss = new ServerSocket(port);
        System.out.println("listening " + ss.getLocalPort());
        while (true) {
            Socket s = ss.accept();
            Thread t = new Conn(s);
            t.start();
        }
    }
}

class Conn extends Thread {
    Socket s;
    Conn(Socket s) { this.s = s; }

    public void run() {
        try {
            InputStream in = s.getInputStream();
            OutputStream out = s.getOutputStream();
            byte[] buf = new byte[2048];
            StringBuilder sb = new StringBuilder();
            boolean headDone = false;
            int headEnd = -1;
            while (!headDone) {
                int n = in.read(buf);
                if (n < 0) break;
                sb.append(new String(buf, 0, n));
                headEnd = sb.toString().indexOf("\r\n\r\n");
                if (headEnd >= 0) headDone = true;
            }
            String resp = "HTTP/1.0 200 OK\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n"
                    + headEnd;
            out.write(resp.getBytes());
            out.close();
        } catch (Exception e) {
            System.out.println("conn error " + e.getMessage());
        }
    }
}
