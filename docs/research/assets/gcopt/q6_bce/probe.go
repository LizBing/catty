// Package q6_bce probes which loop shapes make gc's bounds-check elimination
// (BCE) provably remove slice index checks. The evidence is the deterministic
// output of -d=ssa/check_bce/debug=1 plus the -m escape report; benchmarks
// here only confirm the shapes are within noise of each other.
package q6_bce

// sink and a package-level slice so some cases genuinely defeat BCE (a call
// can reassign the global, so the compiler must re-check bounds).
var sink int64

var global = make([]int64, 256)

//go:noinline
func clobberGlobal() {
	global = append(global, 1)
}

// SumRangeIndex uses the range form, which carries the i < len(s) invariant
// for free.
func SumRangeIndex(s []int64) int64 {
	var t int64
	for i := range s {
		t += s[i]
	}
	return t
}

// SumRangeValue iterates values directly, needing no index check at all.
func SumRangeValue(s []int64) int64 {
	var t int64
	for _, v := range s {
		t += v
	}
	return t
}

// SumForLen reads len(s) in the condition each iteration.
func SumForLen(s []int64) int64 {
	var t int64
	for i := 0; i < len(s); i++ {
		t += s[i]
	}
	return t
}

// SumForCache hoists len(s) into a local.
func SumForCache(s []int64) int64 {
	var t int64
	n := len(s)
	for i := 0; i < n; i++ {
		t += s[i]
	}
	return t
}

// SumReverse iterates downwards from len(s)-1.
func SumReverse(s []int64) int64 {
	var t int64
	for i := len(s) - 1; i >= 0; i-- {
		t += s[i]
	}
	return t
}

// SumGlobalCall indexes a package-level slice with a noinline call inside the
// loop, which can reassign the global and thus defeats BCE on global[i].
func SumGlobalCall() int64 {
	var t int64
	for i := 0; i < len(global); i++ {
		t += global[i]
		clobberGlobal()
	}
	return t
}

// SumGlobalCallBefore puts the mutating call *before* the access, so the
// compiler cannot reuse the loop guard: global may have been reassigned to a
// shorter slice, and global[i] keeps its bounds check.
func SumGlobalCallBefore() int64 {
	var t int64
	for i := 0; i < len(global); i++ {
		clobberGlobal()
		t += global[i]
	}
	return t
}

// SumWithN iterates to an externally supplied n that the compiler cannot
// relate to len(s), so s[i] keeps its bounds check (clean, call-free cost
// probe for the check itself).
func SumWithN(s []int64, n int) int64 {
	var t int64
	for i := 0; i < n; i++ {
		t += s[i]
	}
	return t
}

// SumOpaqueIndex indexes s with a value from a noinline helper, so no proof
// can relate it to len(s).
func SumOpaqueIndex(s []int64) int64 {
	var t int64
	for i := 0; i < len(s); i++ {
		t += s[opaqueIndex(i)]
	}
	return t
}

//go:noinline
func opaqueIndex(i int) int {
	return i & 255
}
