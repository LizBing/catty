package gate5bench;

/** One element of the deterministic application object graph. */
final class Node {
    int value;
    int weight;
    Node next;
    int[] buckets;
}
