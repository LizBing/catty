package gate5bench;

/**
 * Workload entry point. The same class files run under OpenJDK 21 and under
 * every Catty engine configuration; the command argument selects the work and
 * never changes the deterministic operation counts.
 */
public final class Main {
    private Main() {
    }

    public static void main(String[] args) throws Exception {
        String command = args.length > 0 ? args[0] : "verify";
        if (command.equals("verify")) {
            Bench.verify();
        } else if (command.equals("cold")) {
            Bench.cold();
        } else if (command.equals("steady")) {
            Bench.steady(Sizes.FULL_BATCHES, Sizes.FULL_WARMUP);
        } else if (command.equals("steady-smoke")) {
            Bench.steady(Sizes.SMOKE_BATCHES, Sizes.SMOKE_WARMUP);
        } else if (command.equals("rss")) {
            Bench.rss();
        } else if (command.equals("retained")) {
            Retained.run();
        } else {
            throw new IllegalArgumentException("unknown benchmark command");
        }
    }
}
