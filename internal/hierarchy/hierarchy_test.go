package hierarchy

import (
	"errors"
	"fmt"
	"testing"
)

func chainGraph(n int) ([]string, Reader) {
	nodes := map[string]Node{}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("C%03d", i)
		if i == 0 {
			nodes[name] = Node{Super: "java/lang/Object"}
		} else {
			nodes[name] = Node{Super: fmt.Sprintf("C%03d", i-1)}
		}
	}
	nodes["java/lang/Object"] = Node{}
	read := func(name string) (Node, bool) {
		n, ok := nodes[name]
		return n, ok
	}
	roots := []string{fmt.Sprintf("C%03d", n-1)}
	return roots, read
}

func TestCheckExactDepthBoundaryIsOrderIndependent(t *testing.T) {
	for _, n := range []int{128, 129} {
		roots, read := chainGraph(n)
		err := Check(roots, read, 128, nil)
		if n == 128 {
			if err != nil {
				t.Fatalf("128-edge chain rejected: %v", err)
			}
		} else {
			var depth *DepthError
			if !errors.As(err, &depth) || depth.Depth != 129 {
				t.Fatalf("129-edge chain: expected depth 129, got %v", err)
			}
		}
		// Ancestor-first visitation memoizes every ancestor; the reviewer
		// reproducer supplied roots in that order.
		ancestors := make([]string, 0, n)
		for i := 0; i < n; i++ {
			ancestors = append(ancestors, fmt.Sprintf("C%03d", i))
		}
		err2 := Check(ancestors, read, 128, nil)
		if n == 128 {
			if err2 != nil {
				t.Fatalf("ancestor-first 128-edge chain rejected: %v", err2)
			}
		} else {
			var depth *DepthError
			if !errors.As(err2, &depth) || depth.Depth != 129 {
				t.Fatalf("ancestor-first 129-edge chain: expected depth 129, got %v", err2)
			}
		}
	}
}

func TestCheckInterfacePath(t *testing.T) {
	// A class implementing the deepest interface contributes the class edge
	// plus the interface-extension edges.
	nodes := map[string]Node{
		"Impl":   {Interfaces: []string{"I126"}},
		"I000":   {Super: "java/lang/Object"},
		"Object": {},
	}
	for i := 1; i <= 126; i++ {
		nodes[fmt.Sprintf("I%03d", i)] = Node{Super: "java/lang/Object", Interfaces: []string{fmt.Sprintf("I%03d", i-1)}}
	}
	read := func(name string) (Node, bool) { n, ok := nodes[name]; return n, ok }
	if err := Check([]string{"Impl"}, read, 128, nil); err != nil {
		t.Fatalf("128-edge interface path rejected: %v", err)
	}
	nodes["I127"] = Node{Super: "java/lang/Object", Interfaces: []string{"I126"}}
	nodes["Impl"] = Node{Interfaces: []string{"I127"}}
	if err := Check([]string{"Impl"}, read, 128, nil); err == nil {
		t.Fatal("129-edge interface path accepted")
	}
}

func TestCheckCycleAndChargeRefusal(t *testing.T) {
	nodes := map[string]Node{
		"A": {Super: "B"},
		"B": {Super: "A"},
	}
	read := func(name string) (Node, bool) { n, ok := nodes[name]; return n, ok }
	var cycle *CycleError
	if err := Check([]string{"A"}, read, 128, nil); !errors.As(err, &cycle) {
		t.Fatalf("cycle not detected: %v", err)
	}

	calls := 0
	refused := errors.New("refused")
	charge := func() error {
		calls++
		if calls == 2 {
			return refused
		}
		return nil
	}
	if err := Check([]string{"A"}, read, 128, charge); !errors.Is(err, refused) {
		t.Fatalf("charge refusal not propagated: %v", err)
	}
	if calls != 2 {
		t.Fatalf("charge requested %d times after refusal, want 2", calls)
	}
}
