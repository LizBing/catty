package admission

import (
	"errors"
	"strings"

	"catty/internal/classfile"
	"catty/internal/hierarchy"
	"catty/internal/verify"
	"catty/rt"
)

// parentClassTable builds the read-only class metadata view visible from a
// parent loader. Nearest definitions win; the candidate must not shadow any of
// them.
func parentClassTable(parent *rt.Loader) map[string]*classfile.Class {
	out := map[string]*classfile.Class{}
	chain := parent.Chain()
	for i := len(chain) - 1; i >= 0; i-- {
		l := chain[i]
		for name, c := range l.Classes {
			if _, ok := out[name]; !ok {
				out[name] = classfileOf(c)
			}
		}
	}
	return out
}

func classfileOf(c *rt.Class) *classfile.Class {
	out := &classfile.Class{
		Name:        c.Name,
		Flags:       c.Flags,
		NestHost:    c.NestHost,
		NestMembers: append([]string(nil), c.NestMembers...),
		Major:       65,
	}
	if c.Super != nil {
		out.Super = c.Super.Name
	}
	for _, in := range c.Interfaces {
		out.Interfaces = append(out.Interfaces, in.Name)
	}
	for _, f := range c.StaticField {
		out.Fields = append(out.Fields, &classfile.Field{
			Owner: c.Name, Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags, Constant: f.Constant,
		})
	}
	for _, f := range c.InstanceField {
		out.Fields = append(out.Fields, &classfile.Field{
			Owner: c.Name, Name: f.Name, Descriptor: f.Descriptor, Flags: f.Flags,
		})
	}
	for _, m := range c.Methods {
		out.Methods = append(out.Methods, &classfile.Method{
			Owner: c.Name, Name: m.Name, Descriptor: m.Descriptor, Flags: m.Flags,
			MaxStack: m.MaxStack, MaxLocals: m.MaxLocals, CodeLength: m.CodeSize,
			Insns: m.Insns, Refs: m.Refs, Handlers: m.Handlers,
		})
	}
	return out
}

func lookupMeta(name string, byName map[string]*classfile.Class, parent map[string]*classfile.Class) (*classfile.Class, bool) {
	c, _, ok := lookupEffective(name, byName, parent)
	return c, ok
}

// lookupEffective resolves parent-first and reports whether the winning
// definition belongs to the parent chain (and therefore a different defining
// loader) or to the candidate's own loader.
func lookupEffective(name string, byName, parent map[string]*classfile.Class) (*classfile.Class, bool, bool) {
	if c, ok := parent[name]; ok {
		return c, true, true
	}
	if c, ok := byName[name]; ok {
		return c, false, true
	}
	return nil, false, false
}

// candidateAccessible checks class accessibility from a candidate defined in a
// new plugin loader. A package-private target is accessible only when the
// candidate's own loader defines it in the same runtime package.
func candidateAccessible(candidate, target *classfile.Class, targetFromParent bool) bool {
	if target.Flags&rt.AccPublic != 0 {
		return true
	}
	if targetFromParent {
		return false
	}
	return packageOf(candidate.Name) == packageOf(target.Name)
}

func checkCandidate(a *archive, candidate []*sourceClass, parent map[string]*classfile.Class,
	byName map[string]*classfile.Class, profile string, limits Limits, work *verify.Work) *rt.PluginError {
	charge := func(n int64, loc location, feature string) *rt.PluginError {
		if work == nil {
			return nil
		}
		if err := work.Charge(n, feature); err != nil {
			loc.Feature = feature
			loc.Limit = "verificationWork"
			return locFailure("LIMIT", loc, err.Error())
		}
		return nil
	}
	for _, sc := range candidate {
		c := sc.class
		loc := sc.loc
		if perr := charge(1, loc, "class table element"); perr != nil {
			return perr
		}
		for range c.Fields {
			if perr := charge(1, loc, "field table element"); perr != nil {
				return perr
			}
		}
		for _, m := range c.Methods {
			if perr := charge(1, loc, "method table element"); perr != nil {
				return perr
			}
			for range m.Refs {
				if perr := charge(1, loc, "constant-pool reference"); perr != nil {
					return perr
				}
			}
			for range m.Insns {
				if perr := charge(1, loc, "instruction table element"); perr != nil {
					return perr
				}
			}
			for range m.Handlers {
				if perr := charge(1, loc, "handler table element"); perr != nil {
					return perr
				}
			}
		}
		if c.Name == "java/lang/Object" {
			return locFailure("PACKAGE", loc, "candidate redefines java/lang/Object")
		}
		if c.Super != "" {
			super, fromParent, ok := lookupEffective(c.Super, byName, parent)
			if !ok {
				return locFailure("DEPENDENCY", loc, "missing superclass "+c.Super)
			}
			if super.Flags&rt.AccInterface != 0 {
				return locFailure("VERIFY", loc, "superclass is an interface: "+c.Super)
			}
			if super.Flags&rt.AccFinal != 0 {
				return locFailure("VERIFY", loc, "superclass is final: "+c.Super)
			}
			if !candidateAccessible(c, super, fromParent) {
				return locFailure("VERIFY", loc, "superclass is not accessible from the plugin loader: "+c.Super)
			}
		} else if c.Flags&rt.AccInterface == 0 {
			return locFailure("VERIFY", loc, "class has no superclass")
		}
		for _, in := range c.Interfaces {
			iface, fromParent, ok := lookupEffective(in, byName, parent)
			if !ok {
				return locFailure("DEPENDENCY", loc, "missing interface "+in)
			}
			if iface.Flags&rt.AccInterface == 0 {
				return locFailure("VERIFY", loc, "implemented type is not an interface: "+in)
			}
			if !candidateAccessible(c, iface, fromParent) {
				return locFailure("VERIFY", loc, "interface is not accessible from the plugin loader: "+in)
			}
		}
		for _, f := range c.Fields {
			if _, err := rt.ParseFieldType(f.Descriptor); err != nil {
				return locFailure("CLASS_FORMAT", loc, "field "+f.Name+": "+err.Error())
			}
		}
		for _, m := range c.Methods {
			mloc := loc
			mloc.Method, mloc.Descriptor, mloc.PC = m.Name, m.Descriptor, 0
			if m.Flags&rt.AccSynchronized != 0 {
				return locFailure("PROFILE", mloc, "synchronized methods are unsupported")
			}
			if m.Flags&rt.AccNative != 0 {
				return locFailure("PROFILE", mloc, "native methods are unsupported")
			}
			if _, err := rt.ParseMethodType(m.Descriptor); err != nil {
				return locFailure("CLASS_FORMAT", mloc, err.Error())
			}
		}
	}
	return nil
}

// checkHierarchy applies the effective hierarchy rule to the classes the
// candidate loader can actually define. The reader follows parent-loader
// definitions, so a path that spans a published parent loader counts every
// class/interface edge, and the memoized longest-path traversal is independent
// of archive entry order. Work is charged before each newly visited node.
func checkHierarchy(effective []*classfile.Class, byName, parent map[string]*classfile.Class,
	limits Limits, work *verify.Work) *rt.PluginError {
	roots := make([]string, 0, len(effective))
	for _, c := range effective {
		roots = append(roots, c.Name)
	}
	charge := func() error {
		if work == nil {
			return nil
		}
		return work.Charge(1, "hierarchy visit")
	}
	read := func(name string) (hierarchy.Node, bool) {
		c, _, ok := lookupEffective(name, byName, parent)
		if !ok {
			return hierarchy.Node{}, false
		}
		return hierarchy.Node{Super: c.Super, Interfaces: c.Interfaces}, true
	}
	err := hierarchy.Check(roots, read, limits.HierarchyDepth, charge)
	if err == nil {
		return nil
	}
	var depthErr *hierarchy.DepthError
	if errors.As(err, &depthErr) {
		loc := location{Class: depthErr.Class, Feature: "hierarchy depth", Limit: "hierarchyDepth"}
		return locFailure("LIMIT", loc, depthErr.Error())
	}
	var cycleErr *hierarchy.CycleError
	if errors.As(err, &cycleErr) {
		return failure("VERIFY", "", cycleErr.Class, cycleErr.Error())
	}
	var workErr *verify.WorkError
	if errors.As(err, &workErr) {
		loc := location{Limit: "verificationWork"}
		return locFailure("LIMIT", loc, err.Error())
	}
	return failure("VERIFY", "", "", err.Error())
}

func packageOf(name string) string {
	if i := strings.LastIndexByte(name, '/'); i >= 0 {
		return name[:i]
	}
	return ""
}

func checkNests(effective []*classfile.Class, byName map[string]*classfile.Class,
	parent map[string]*classfile.Class, l *rt.Loader, maxDepth int) *rt.PluginError {
	for _, c := range effective {
		loc := location{Class: c.Name}
		hostName := c.NestHost
		if hostName == "" && len(c.NestMembers) == 0 {
			continue
		}
		if hostName != "" {
			if _, ok := byName[hostName]; ok {
				host := byName[hostName]
				if !containsName(host.NestMembers, c.Name) {
					return locFailure("VERIFY", loc, "nonreciprocal nest membership: "+c.Name+" claims "+hostName)
				}
			} else if _, ok := parent[hostName]; ok {
				// A plugin class must not join a parent-loader nest.
				return locFailure("VERIFY", loc, "cross-loader nest claim to "+hostName)
			} else {
				return locFailure("DEPENDENCY", loc, "missing nest host "+hostName)
			}
		}
		for _, member := range c.NestMembers {
			m, ok := byName[member]
			if !ok {
				if _, parentHas := parent[member]; parentHas {
					return locFailure("VERIFY", loc, "cross-loader nest member "+member)
				}
				return locFailure("DEPENDENCY", loc, "missing nest member "+member)
			}
			if m.NestHost != c.Name {
				return locFailure("VERIFY", loc, "nonreciprocal nest member "+member)
			}
		}
	}
	return nil
}

func containsName(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}

func checkBootstrapReferences(effective []*classfile.Class, byName, parent map[string]*classfile.Class, profile string) *rt.PluginError {
	for _, c := range effective {
		for _, m := range c.Methods {
			for _, ref := range m.Refs {
				loc := location{Class: c.Name, Method: m.Name, Descriptor: m.Descriptor}
				switch ref.Kind {
				case rt.RefClass:
					if reservedDefinition(ref.Class) {
						if _, ok := lookupMeta(ref.Class, byName, parent); !ok {
							loc.Feature = ref.Class
							return locFailure("PROFILE", loc, "unavailable bootstrap class "+ref.Class)
						}
					}
				case rt.RefField, rt.RefMethod:
					if reservedDefinition(ref.Owner) {
						owner, ok := lookupMeta(ref.Owner, byName, parent)
						if !ok {
							loc.Feature = ref.Owner
							return locFailure("PROFILE", loc, "unavailable bootstrap owner "+ref.Owner)
						}
						if ref.Kind == rt.RefField {
							if !hasField(owner, ref.Name, ref.Descriptor, byName, parent) {
								loc.Feature = ref.Owner + "." + ref.Name
								return locFailure("PROFILE", loc, "unavailable bootstrap field")
							}
						} else {
							if !hasMethod(owner, ref.Name, ref.Descriptor, byName, parent) {
								loc.Feature = ref.Owner + "." + ref.Name + ref.Descriptor
								return locFailure("PROFILE", loc, "unavailable bootstrap method")
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// resolveMemberFlags finds a method or field declaration for the protected
// receiver check. Methods and fields are both searched through the hierarchy.
func resolveMemberFlags(owner, name, descriptor string, byName, parent map[string]*classfile.Class) (uint16, string, bool) {
	seen := map[string]bool{}
	var visit func(nameOfClass string) (uint16, string, bool)
	visit = func(nameOfClass string) (uint16, string, bool) {
		if nameOfClass == "" || seen[nameOfClass] {
			return 0, "", false
		}
		seen[nameOfClass] = true
		c, ok := lookupMeta(nameOfClass, byName, parent)
		if !ok {
			return 0, "", false
		}
		for _, m := range c.Methods {
			if m.Name == name && m.Descriptor == descriptor {
				return m.Flags, c.Name, true
			}
		}
		for _, f := range c.Fields {
			if f.Name == name && f.Descriptor == descriptor {
				return f.Flags, c.Name, true
			}
		}
		for _, in := range c.Interfaces {
			if flags, decl, ok := visit(in); ok {
				return flags, decl, true
			}
		}
		return visit(c.Super)
	}
	return visit(owner)
}

func hasField(c *classfile.Class, name, descriptor string, byName, parent map[string]*classfile.Class) bool {
	seen := map[string]bool{}
	var visit func(cur *classfile.Class) bool
	visit = func(cur *classfile.Class) bool {
		if cur == nil || seen[cur.Name] {
			return false
		}
		seen[cur.Name] = true
		for _, f := range cur.Fields {
			if f.Name == name && f.Descriptor == descriptor {
				return true
			}
		}
		for _, in := range cur.Interfaces {
			if cand, ok := lookupMeta(in, byName, parent); ok && visit(cand) {
				return true
			}
		}
		if cur.Super == "" {
			return false
		}
		sup, ok := lookupMeta(cur.Super, byName, parent)
		return ok && visit(sup)
	}
	return visit(c)
}

func hasMethod(c *classfile.Class, name, descriptor string, byName, parent map[string]*classfile.Class) bool {
	seen := map[string]bool{}
	var visit func(cur *classfile.Class) bool
	visit = func(cur *classfile.Class) bool {
		if cur == nil || seen[cur.Name] {
			return false
		}
		seen[cur.Name] = true
		for _, m := range cur.Methods {
			if m.Name == name && m.Descriptor == descriptor {
				return true
			}
		}
		for _, in := range cur.Interfaces {
			if cand, ok := lookupMeta(in, byName, parent); ok && visit(cand) {
				return true
			}
		}
		if cur.Super == "" {
			return false
		}
		sup, ok := lookupMeta(cur.Super, byName, parent)
		return ok && visit(sup)
	}
	return visit(c)
}
