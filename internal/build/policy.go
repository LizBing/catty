package build

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// MethodKey identifies one application method selection target.
type MethodKey struct {
	Owner      string
	Name       string
	Descriptor string
}

// Policy is an ordered set of forced per-method engines.
type Policy struct {
	order      []MethodKey
	selections map[MethodKey]string
}

type policyFile struct {
	Methods []policyEntry `json:"methods"`
}

type policyEntry struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Descriptor string `json:"descriptor"`
	Engine     string `json:"engine"`
}

// LoadPolicy parses and validates the policy JSON shape. Target existence and
// feasibility are checked against the verified class set by ApplyPolicy.
func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("catty: cannot read policy %s: %w", path, err)
	}
	var file policyFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("catty: malformed policy %s: %w", path, err)
	}
	p := &Policy{selections: map[MethodKey]string{}}
	for _, e := range file.Methods {
		key := MethodKey{Owner: e.Owner, Name: e.Name, Descriptor: e.Descriptor}
		if e.Owner == "" || e.Name == "" || e.Descriptor == "" {
			return nil, fmt.Errorf("catty: policy entry is missing owner, name or descriptor")
		}
		if e.Engine != "aot" && e.Engine != "interp" {
			return nil, fmt.Errorf("catty: policy for %s.%s%s has unsupported engine %q",
				e.Owner, e.Name, e.Descriptor, e.Engine)
		}
		if old, ok := p.selections[key]; ok {
			if old == e.Engine {
				return nil, fmt.Errorf("catty: duplicate method selection for %s.%s%s",
					e.Owner, e.Name, e.Descriptor)
			}
			return nil, fmt.Errorf("catty: conflicting selections for %s.%s%s (%s versus %s)",
				e.Owner, e.Name, e.Descriptor, old, e.Engine)
		}
		p.selections[key] = e.Engine
		p.order = append(p.order, key)
	}
	return p, nil
}

// Keys returns the selection keys in policy order.
func (p *Policy) Keys() []MethodKey {
	if p == nil {
		return nil
	}
	out := append([]MethodKey(nil), p.order...)
	return out
}

// Engine returns the forced engine for key.
func (p *Policy) Engine(key MethodKey) (string, bool) {
	if p == nil {
		return "", false
	}
	e, ok := p.selections[key]
	return e, ok
}

// SortedKeys is a deterministic helper for diagnostics.
func (p *Policy) SortedKeys() []MethodKey {
	keys := p.Keys()
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Owner != keys[j].Owner {
			return keys[i].Owner < keys[j].Owner
		}
		if keys[i].Name != keys[j].Name {
			return keys[i].Name < keys[j].Name
		}
		return keys[i].Descriptor < keys[j].Descriptor
	})
	return keys
}
