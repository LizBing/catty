package gate5bench;

/**
 * Machine-readable benchmark output.
 *
 * <p>Catty's Java 21 subset has no string concatenation (it would need
 * {@code invokedynamic}) and no {@code StringBuilder}, so every logical record
 * is emitted as a tag line followed by its values, one value per line. The
 * Python runner in {@code tests/gate5_bench.py} parses exactly this sequence.
 * Numbers are printed with the available {@code PrintStream.println(int)} and
 * {@code println(long)} overloads.
 */
final class Report {
    private Report() {
    }

    static void begin() {
        System.out.println("CATTY-GATE5-BENCH");
        System.out.println(1);
    }

    static void end() {
        System.out.println("END");
        System.out.println(0);
    }

    static void word(String text) {
        System.out.println(text);
    }

    static void number(int value) {
        System.out.println(value);
    }

    static void bignumber(long value) {
        System.out.println(value);
    }
}
