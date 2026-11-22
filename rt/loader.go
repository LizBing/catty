package rt

// Lookup resolves name through the loader chain, parent first. It performs no
// array or primitive handling; classForLoaderName owns those cases.
func (l *Loader) Lookup(name string) (*Class, bool) {
	for cur := l; cur != nil; cur = cur.Parent {
		if c, ok := cur.Classes[name]; ok {
			return c, true
		}
	}
	return nil, false
}

// Local returns a class defined by this loader only.
func (l *Loader) Local(name string) (*Class, bool) {
	if l == nil {
		return nil, false
	}
	c, ok := l.Classes[name]
	return c, ok
}

// Chain returns this loader followed by its ancestors, nearest first.
func (l *Loader) Chain() []*Loader {
	var out []*Loader
	for cur := l; cur != nil; cur = cur.Parent {
		out = append(out, cur)
	}
	return out
}

// IsAncestor reports whether ancestor is l or one of its parents.
func (l *Loader) IsAncestor(ancestor *Loader) bool {
	for cur := l; cur != nil; cur = cur.Parent {
		if cur == ancestor {
			return true
		}
	}
	return false
}

// RuntimeLoader returns the owning runtime for a loader handle.
func (l *Loader) RuntimeLoader() *Runtime {
	if l == nil {
		return nil
	}
	return l.Runtime
}
