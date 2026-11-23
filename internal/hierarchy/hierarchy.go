// Package hierarchy owns the effective superclass/interface graph rule shared
// by image admission, plugin-staging verification and runtime revalidation.
package hierarchy

import "fmt"

// Node is the hierarchy metadata of one class or interface. A missing Super is
// a terminal node.
type Node struct {
	Super      string
	Interfaces []string
}

// Reader returns the metadata for one name. ok=false means the graph ends at
// that name, which is how profiles describe bootstrap roots that have no
// further superclass or interface edges.
type Reader func(name string) (Node, bool)

// DepthError reports a longest path that exceeds the inclusive profile limit.
type DepthError struct {
	Class string
	Depth int
	Limit int
}

func (e *DepthError) Error() string {
	return fmt.Sprintf("hierarchy depth %d exceeds %d at %s", e.Depth, e.Limit, e.Class)
}

// CycleError reports a superclass/interface cycle.
type CycleError struct {
	Class string
}

func (e *CycleError) Error() string {
	return "class hierarchy cycle at " + e.Class
}

// Check visits the graph reachable from roots, charges one unit per newly
// visited node before reading it, rejects cycles, and returns a DepthError when
// the longest path from a root contains more than limit edges.
//
// The traversal memoizes the longest path from every completed node, so an
// ancestor reached through several descendants is evaluated once and its full
// remaining length is still counted for each path. Visit order therefore does
// not affect the result.
func Check(roots []string, read Reader, limit int, charge func() error) error {
	if limit < 0 {
		limit = 0
	}
	const (
		unvisited = 0
		visiting  = 1
		done      = 2
	)
	state := map[string]uint8{}
	longest := map[string]int{}

	var visit func(name string) (int, error)
	visit = func(name string) (int, error) {
		switch state[name] {
		case visiting:
			return 0, &CycleError{Class: name}
		case done:
			return longest[name], nil
		}
		if charge != nil {
			if err := charge(); err != nil {
				return 0, err
			}
		}
		state[name] = visiting
		node, ok := read(name)
		max := 0
		if ok {
			if node.Super != "" {
				depth, err := visit(node.Super)
				if err != nil {
					return 0, err
				}
				if depth+1 > max {
					max = depth + 1
				}
			}
			for _, in := range node.Interfaces {
				depth, err := visit(in)
				if err != nil {
					return 0, err
				}
				if depth+1 > max {
					max = depth + 1
				}
			}
		}
		state[name] = done
		longest[name] = max
		if max > limit {
			return 0, &DepthError{Class: name, Depth: max, Limit: limit}
		}
		return max, nil
	}
	for _, root := range roots {
		if _, err := visit(root); err != nil {
			return err
		}
	}
	return nil
}
