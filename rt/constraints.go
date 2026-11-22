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

// conflictWith checks the constraints involving key against one candidate
// initiation map. It reports the first peer whose resolved class differs.
func conflictWith(constraints map[constraintKey]struct{}, initiations map[initiationKey]*Class,
	key initiationKey, class *Class) bool {
	for ck := range constraints {
		if ck.A.Name != key.Name {
			continue
		}
		var peer initiationKey
		switch key {
		case ck.A:
			peer = ck.B
		case ck.B:
			peer = ck.A
		default:
			continue
		}
		if peer.Loader == key.Loader {
			continue
		}
		if peerClass, ok := initiations[peer]; ok && peerClass != class {
			return true
		}
	}
	return false
}

// recordInitiation commits one successful loader resolution after checking all
// equality constraints already registered for the name. A conflict leaves both
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
	if conflictWith(r.constraints, r.initiations, key, c) {
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

// stageConstraints resolves each leaf through both loaders on copies of the
// constraint maps. It publishes initiation records and constraints only when
// no conflict was observed. Absent classes keep ordinary resolution timing and
// are not turned into a constraint failure here.
func (r *Runtime) stageConstraints(a, b *Loader, types []Type) bool {
	if a == nil || b == nil || a == b {
		return false
	}
	arrays := make(map[*Class]*Class, len(r.arrayClasses))
	for k, v := range r.arrayClasses {
		arrays[k] = v
	}
	primitiveArrays := make(map[byte]*Class, len(r.primitiveArrayClass))
	for k, v := range r.primitiveArrayClass {
		primitiveArrays[k] = v
	}
	identity := r.nextIdentity
	conflict := r.stageConstraintsTx(a, b, types)
	if conflict {
		r.arrayClasses = arrays
		r.primitiveArrayClass = primitiveArrays
		r.nextIdentity = identity
	}
	return conflict
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
		ca, errA := r.classForLoaderName(a, name)
		if errA != nil {
			continue
		}
		cb, errB := r.classForLoaderName(b, name)
		if errB != nil {
			continue
		}
		if ca != cb {
			return true
		}
		keyA := initiationKey{Loader: a, Name: name}
		keyB := initiationKey{Loader: b, Name: name}
		if existing, ok := initiations[keyA]; ok && existing != ca {
			return true
		}
		if existing, ok := initiations[keyB]; ok && existing != cb {
			return true
		}
		initiations[keyA] = ca
		initiations[keyB] = cb
		ck := makeConstraintKey(keyA, keyB)
		constraints[ck] = struct{}{}
		if conflictWith(constraints, initiations, keyA, ca) ||
			conflictWith(constraints, initiations, keyB, cb) {
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
