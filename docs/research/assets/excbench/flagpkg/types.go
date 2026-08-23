// Package flagpkg implements candidate F: every method that can throw returns
// a (result, *common.JException) pair, and every call site checks the flag
// explicitly. A throw is a plain return of the sentinel; propagation is an
// if err != nil { return err } at each boundary.
package flagpkg

import "excbench/common"

// Handler results, chosen distinct so a correctness test can tell the layers
// apart; values mirror panicpkg for comparability.
const (
	resultA = -200
	resultC = -100
)

// --- Normal path: d never errs --------------------------------------------

//go:noinline
func aNorm(x int) (int, *common.JException) {
	v, err := bNorm(x)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//go:noinline
func bNorm(x int) (int, *common.JException) {
	v, err := cNorm(x)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//go:noinline
func cNorm(x int) (int, *common.JException) {
	v, err := dNorm(x)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//go:noinline
func dNorm(x int) (int, *common.JException) {
	return common.Work(x), nil
}

// --- Shallow throw-catch: d errs, c catches (one boundary) -----------------

//go:noinline
func dShallow(x int) (int, *common.JException) {
	return 0, common.SentinelC
}

//go:noinline
func cShallow(x int) (int, *common.JException) {
	v, err := dShallow(x)
	if err != nil {
		return resultC, nil
	}
	return v, nil
}

// --- Deep propagation: d errs, caught at a --------------------------------

//go:noinline
func aDeep(x int) (int, *common.JException) {
	v, err := bDeep(x)
	if err != nil {
		return resultA, nil
	}
	return v, nil
}

//go:noinline
func bDeep(x int) (int, *common.JException) {
	v, err := cDeep(x)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//go:noinline
func cDeep(x int) (int, *common.JException) {
	v, err := dDeep(x)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//go:noinline
func dDeep(x int) (int, *common.JException) {
	return 0, common.SentinelA
}
