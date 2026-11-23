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

// constraintTx stages initiation and equality updates on detached copies. A
// refused staging operation leaves the transaction itself unchanged, and only
// commit publishes staged changes to the runtime. Callers discard the
// transaction when the surrounding operation fails, so a rejected plugin
// candidate cannot leave pending constraints behind.
type constraintTx struct {
	initiations map[initiationKey]*Class
	constraints map[constraintKey]struct{}
}

func (r *Runtime) newConstraintTx() *constraintTx {
	r.ensureConstraintState()
	tx := &constraintTx{
		initiations: make(map[initiationKey]*Class, len(r.initiations)),
		constraints: make(map[constraintKey]struct{}, len(r.constraints)),
	}
	for k, v := range r.initiations {
		tx.initiations[k] = v
	}
	for k := range r.constraints {
		tx.constraints[k] = struct{}{}
	}
	return tx
}

// commit hands the staged maps to the runtime. The transaction must not be used
// again afterwards.
func (tx *constraintTx) commit(r *Runtime) {
	r.initiations = tx.initiations
	r.constraints = tx.constraints
}

// constraintConflict reports two distinct resolved classes that a loading
// constraint would require to be the same class.
type constraintConflict struct {
	Name string
	A, B *Class
}

// componentConflict walks the full equality component that contains key and
// reports the first pair of resolved peers that disagree. pending supplies the
// prospective resolution of key; nil means the recorded initiation is used.
// Unresolved peers keep ordinary lazy resolution timing and impose no immediate
// failure.
func componentConflict(constraints map[constraintKey]struct{}, initiations map[initiationKey]*Class,
	key initiationKey, pending *Class) *constraintConflict {
	var first *Class
	seen := map[initiationKey]bool{key: true}
	queue := []initiationKey{key}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		c := initiations[cur]
		if cur == key && pending != nil {
			c = pending
		}
		if c != nil {
			if first == nil {
				first = c
			} else if first != c {
				return &constraintConflict{Name: key.Name, A: first, B: c}
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
	return nil
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
	if componentConflict(r.constraints, r.initiations, key, c) != nil {
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

// stageConstraints adds an equality edge for each reference leaf type. It never
// resolves a descriptor leaf just to decide whether the edge can be staged: when
// either side has not initiated the leaf yet the equality stays pending, and
// recordInitiation enforces it when a later lookup publishes one side. Primitive
// and void descriptors contribute no edge. All updates are built on copies and
// published together only on success.
func (r *Runtime) stageConstraints(a, b *Loader, types []Type) bool {
	if a == nil || b == nil || a == b {
		return false
	}
	return r.stageConstraintsTx(a, b, types)
}

func (r *Runtime) stageConstraintsTx(a, b *Loader, types []Type) bool {
	tx := r.newConstraintTx()
	if tx.stage(a, b, types) != nil {
		return true
	}
	tx.commit(r)
	return false
}

// stage adds one equality edge for every reference leaf type and validates the
// merged equality component as a whole. Joining two previously disjoint groups
// is refused when their already-resolved members disagree, even when both new
// endpoints are still unresolved. The transaction is unchanged on conflict, so
// an earlier edge from the same descriptor never leaks out of a refused
// operation.
func (tx *constraintTx) stage(a, b *Loader, types []Type) *constraintConflict {
	if a == nil || b == nil || a == b {
		return nil
	}
	next := make(map[constraintKey]struct{}, len(tx.constraints)+len(types))
	for k := range tx.constraints {
		next[k] = struct{}{}
	}
	for _, t := range types {
		name := leafClassName(t)
		if name == "" {
			continue
		}
		keyA := initiationKey{Loader: a, Name: name}
		keyB := initiationKey{Loader: b, Name: name}
		next[makeConstraintKey(keyA, keyB)] = struct{}{}
		if conflict := componentConflict(next, tx.initiations, keyA, nil); conflict != nil {
			return conflict
		}
	}
	tx.constraints = next
	return nil
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
