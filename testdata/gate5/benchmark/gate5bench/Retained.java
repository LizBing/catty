package gate5bench;

/**
 * Retained-graph workload for the explicit GC observation. The graph is held by
 * a static root for the whole measurement; the checksum is computed before and
 * after the explicit collection so a collected or corrupted graph is visible in
 * the workload output as well as in the runtime heap observation.
 */
final class Retained {
    private static Node[] roots;
    private static int retainedCount;

    private Retained() {
    }

    static int checksum(Node[] nodes) {
        int acc = 0;
        int i = 0;
        while (i < nodes.length) {
            Node node = nodes[i];
            acc = acc * 31 + node.value + node.weight;
            int j = 0;
            while (j < 4) {
                acc = acc * 17 + node.buckets[j];
                j = j + 1;
            }
            i = i + 1;
        }
        return acc;
    }

    static void run() {
        int count = Sizes.RETAINED_NODES;
        Node head = Workload.buildGraph(count, Sizes.SEED);
        Node[] graph = new Node[count];
        Node node = head;
        int i = 0;
        while (i < count) {
            graph[i] = node;
            node = node.next;
            i = i + 1;
        }
        roots = graph;
        retainedCount = count;
        int before = checksum(graph);
        Runtime.getRuntime().gc();
        int after = checksum(roots);
        if (before != after) {
            throw new IllegalStateException("retained graph changed during collection");
        }
        Report.begin();
        Report.word("COMMAND");
        Report.word("retained");
        Report.word("RETAINED");
        Report.number(retainedCount);
        Report.word("CHECK-BEFORE");
        Report.number(before);
        Report.word("GC-CALLS");
        Report.number(1);
        Report.word("CHECK-AFTER");
        Report.number(after);
        Report.end();
    }
}
