package rt

import "strings"

// initiationKey records that loader resolved name to a class. The runtime uses
// this with equality constraints to implement JVMS 5.3.4 loading constraints.
type initiationKey struct {
	Loader *Loader
	Name   string
}

// constraintKey is an unordered pair of loader-specific resolutions of one
// class name. The pair is canonicalized by loader ID.
type constraintKey struct {
	A, B initiationKey
}

func (r *Runtime) ensureConstraintState() {
	if r.initiations == nil {
		r.initiations = map[initiationKey]*Class{}
	}
	if r.constraints == nil {
		r.constraints = map[constraintKey]struct{}{}
	}
}

func makeConstraintKey(a, b initiationKey) constraintKey {
	if a.Loader.ID > b.Loader.ID {
		a, b = b, a
	}
	return constraintKey{A: a, B: b}
}

// connectedConflict walks the full equality component for key/name and
// compares every already-resolved peer with class. A constraint edge may have
// been staged before either endpoint resolved, so a peer reached only through
// another pending endpoint still participates; unresolved peers keep ordinary
// lazy resolution timing and impose no immediate failure.
func connectedConflict(constraints map[constraintKey]struct{}, initiations map[initiationKey]*Class,
	key initiationKey, class *Class) bool {
	if class == nil {
		return false
	}
	seen := map[initiationKey]bool{key: true}
	queue := []initiationKey{key}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur != key {
			if peerClass, ok := initiations[cur]; ok && peerClass != class {
				return true
			}
		}
		for ck := range constraints {
			if ck.A.Name != key.Name {
				continue
			}
			var next initiationKey
			switch cur {
			case ck.A:
				next = ck.B
			case ck.B:
				next = ck.A
			default:
				continue
			}
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return false
}

// recordInitiation commits one successful loader resolution after checking the
// full connected equality component for the binary name. A conflict leaves both
// maps unchanged.
func (r *Runtime) recordInitiation(loader *Loader, name string, c *Class) bool {
	if loader == nil || c == nil {
		return false
	}
	r.ensureConstraintState()
	key := initiationKey{Loader: loader, Name: name}
	if existing, ok := r.initiations[key]; ok {
		return existing != c
	}
	if connectedConflict(r.constraints, r.initiations, key, c) {
		return true
	}
	r.initiations[key] = c
	return false
}

// recordClassInitiation records the class name and, for an array descriptor,
// its reference leaf. It returns true on a constraint conflict.
func (r *Runtime) recordClassInitiation(loader *Loader, name string, c *Class) bool {
	if r.recordInitiation(loader, name, c) {
		return true
	}
	if leaf := arrayLeafName(name); leaf != "" {
		lc, err := r.classForLoaderName(loader, leaf)
		if err == nil && r.recordInitiation(loader, leaf, lc) {
			return true
		}
	}
	return false
}

// imposeMethodConstraints stages equality constraints for every reference leaf
// type in a method descriptor and commits them only if no conflict occurs.
func (r *Runtime) imposeMethodConstraints(callerLoader, declaringLoader *Loader, mt *MethodType) bool {
	if mt == nil || callerLoader == nil || declaringLoader == nil || callerLoader == declaringLoader {
		return false
	}
	types := append([]Type(nil), mt.Params...)
	types = append(types, mt.Return)
	return r.stageConstraints(callerLoader, declaringLoader, types)
}

// imposeFieldConstraints is the field-descriptor equivalent.
func (r *Runtime) imposeFieldConstraints(callerLoader, declaringLoader *Loader, ft Type) bool {
	if callerLoader == nil || declaringLoader == nil || callerLoader == declaringLoader {
		return false
	}
	return r.stageConstraints(callerLoader, declaringLoader, []Type{ft})
}

// stageConstraints adds an equality edge for each reference leaf type. It
// never resolves a descriptor leaf just to decide whether the edge can be
// staged: when either side has not initiated the leaf yet the equality stays
// pending, and recordInitiation enforces it when a later lookup publishes one
// side. Primitive and void descriptors contribute no edge. All updates are
// built on copies and published together only on success.
func (r *Runtime) stageConstraints(a, b *Loader, types []Type) bool {
	if a == nil || b == nil || a == b {
		return false
	}
	return r.stageConstraintsTx(a, b, types)
}

func (r *Runtime) stageConstraintsTx(a, b *Loader, types []Type) bool {
	r.ensureConstraintState()
	initiations := make(map[initiationKey]*Class, len(r.initiations)+len(types)*2)
	for k, v := range r.initiations {
		initiations[k] = v
	}
	constraints := make(map[constraintKey]struct{}, len(r.constraints)+len(types))
	for k := range r.constraints {
		constraints[k] = struct{}{}
	}
	for _, t := range types {
		name := leafClassName(t)
		if name == "" {
			continue
		}
		keyA := initiationKey{Loader: a, Name: name}
		keyB := initiationKey{Loader: b, Name: name}
		ca, okA := initiations[keyA]
		cb, okB := initiations[keyB]
		if okA && okB && ca != cb {
			return true
		}
		if okA && connectedConflict(constraints, initiations, keyA, ca) {
			return true
		}
		if okB && connectedConflict(constraints, initiations, keyB, cb) {
			return true
		}
		constraints[makeConstraintKey(keyA, keyB)] = struct{}{}
		// The new edge can connect two previously disjoint components, so
		// recheck each already-resolved endpoint against the joined component.
		if okA && connectedConflict(constraints, initiations, keyA, ca) {
			return true
		}
		if okB && connectedConflict(constraints, initiations, keyB, cb) {
			return true
		}
	}
	r.initiations = initiations
	r.constraints = constraints
	return false
}

// leafClassName returns the class leaf of a reference descriptor, including
// nested array components. Primitive and void leaves return "".
func arrayLeafName(desc string) string {
	for strings.HasPrefix(desc, "[") {
		desc = desc[1:]
	}
	if strings.HasPrefix(desc, "L") && strings.HasSuffix(desc, ";") {
		return desc[1 : len(desc)-1]
	}
	return ""
}

func leafClassName(t Type) string {
	switch t.Kind {
	case 'L':
		return t.ClassName
	case 'A':
		desc := t.Descriptor
		for strings.HasPrefix(desc, "[") {
			desc = desc[1:]
		}
		if strings.HasPrefix(desc, "L") && strings.HasSuffix(desc, ";") {
			return desc[1 : len(desc)-1]
		}
	}
	return ""
}
