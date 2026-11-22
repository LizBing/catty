package hidden;

// Package-private class with a public method: getMethod finds it, but
// Method.invoke from another package must fail the access check.
class Hidden {
    public int value() {
        return 42;
    }
}
