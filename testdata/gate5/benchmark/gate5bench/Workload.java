package gate5bench;

/**
 * Deterministic object-graph workload: field, array and allocation work over a
 * linked graph, with an engine crossing through {@link Hot#mix} and an AOT
 * callback for every node visit.
 */
final class Workload {
    private Workload() {
    }

    /** Builds the same graph for the same seed and node count on every runtime. */
    static Node buildGraph(int nodes, int seed) {
        Node head = null;
        Node tail = null;
        int state = seed | 1;
        int i = 0;
        while (i < nodes) {
            Node node = new Node();
            state = state * 1664525 + 1013904223;
            node.value = state;
            node.weight = (state >>> 8) & 1023;
            node.buckets = new int[4];
            int local = state;
            int j = 0;
            while (j < 4) {
                local = local * 1103515245 + 12345;
                node.buckets[j] = local;
                j = j + 1;
            }
            if (head == null) {
                head = node;
            } else {
                tail.next = node;
            }
            tail = node;
            i = i + 1;
        }
        return head;
    }

    /**
     * Visits the graph once. This is the AOT caller of the interpreted
     * {@link Hot#mix}; each visit also allocates a small transient array so the
     * graph path includes escaping allocation work.
     */
    public static int process(Node head, int nodes, int rounds) {
        int acc = 17;
        int round = 0;
        while (round < rounds) {
            Node node = head;
            int i = 0;
            while (i < nodes) {
                int[] buckets = node.buckets;
                int b = buckets[i & 3];
                int v = node.value + b + node.weight;
                int[] scratch = new int[2];
                scratch[0] = v;
                scratch[1] = i;
                acc = Hot.mix(acc, v + scratch[1]);
                node.value = node.value + (acc & 7);
                node = node.next;
                i = i + 1;
            }
            round = round + 1;
        }
        return acc;
    }

    /** Deterministic checksum over the whole graph. */
    static int checksum(Node head, int nodes) {
        int acc = 0;
        Node node = head;
        int i = 0;
        while (i < nodes) {
            acc = acc * 31 + node.value + node.weight;
            int j = 0;
            while (j < 4) {
                acc = acc * 17 + node.buckets[j];
                j = j + 1;
            }
            node = node.next;
            i = i + 1;
        }
        return acc;
    }
}
