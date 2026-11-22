package ext;

public class Worker {
    public static int accept(Token token) { return 1; }

    public static int accept(Token[] tokens) { return 2; }

    public static Token[] tokens() { return new Token[1]; }
}
