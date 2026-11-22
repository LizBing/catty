package rt

import "fmt"

// PluginLimits are the published Gate 3 admission ceilings. Tests may inject
// smaller limits through the admission package; production values are profile
// constants and travel in the image report.
type PluginLimits struct {
	CompressedBytes      int64 `json:"compressedBytes"`
	ExpandedBytes        int64 `json:"expandedBytes"`
	EntryBytes           int64 `json:"entryBytes"`
	Entries              int   `json:"entries"`
	ClassEntries         int   `json:"classEntries"`
	HierarchyDepth       int   `json:"hierarchyDepth"`
	VerificationWork     int64 `json:"verificationWork"`
	Loaders              int   `json:"loaders"`
	CumulativeClassBytes int64 `json:"cumulativeClassBytes"`
}

// Profile describes one published Catty runtime profile. The build selects a
// profile, the image carries its identity and ABI, and Install accepts only the
// matching pair before any Java code runs.
type Profile struct {
	ID      string
	ABI     int
	Plugins bool
	Limits  *PluginLimits
}

// Profile identity constants. Gate 2 remains the default so builds without the
// plugin opt-in keep their published contract.
const (
	Gate2ProfileID = "catty-jvm21-gate2"
	Gate2ABI       = 2
	Gate3ProfileID = "catty-jvm21-gate3"
	Gate3ABI       = 3

	// ABIRevision and ProfileID name the default (Gate 2) pair.
	ABIRevision = Gate2ABI
	ProfileID   = Gate2ProfileID
)

// Gate3PluginLimits are the initial profile ceilings from GATE3.md.
var Gate3PluginLimits = &PluginLimits{
	CompressedBytes:      16 << 20,
	ExpandedBytes:        64 << 20,
	EntryBytes:           4 << 20,
	Entries:              4096,
	ClassEntries:         1024,
	HierarchyDepth:       128,
	VerificationWork:     16777216,
	Loaders:              64,
	CumulativeClassBytes: 128 << 20,
}

var (
	Gate2Profile = &Profile{ID: Gate2ProfileID, ABI: Gate2ABI}
	Gate3Profile = &Profile{ID: Gate3ProfileID, ABI: Gate3ABI, Plugins: true, Limits: Gate3PluginLimits}
)

// ProfileByID returns the published profile with the exact identity.
func ProfileByID(id string) (*Profile, bool) {
	switch id {
	case Gate2ProfileID:
		return Gate2Profile, true
	case Gate3ProfileID:
		return Gate3Profile, true
	}
	return nil, false
}

// ProfileForImage validates the profile/ABI pair carried by an image.
func ProfileForImage(profile string, abi int) (*Profile, error) {
	p, ok := ProfileByID(profile)
	if !ok {
		return nil, fmt.Errorf("catty: unsupported image profile %q", profile)
	}
	if p.ABI != abi {
		return nil, fmt.Errorf("catty: image profile %q requires ABI %d, got %d", profile, p.ABI, abi)
	}
	return p, nil
}

// Bootstrap inventory accessors scoped to a profile.

// BootstrapHasClassFor reports whether name belongs to the profile inventory.
func BootstrapHasClassFor(profile, name string) bool {
	p, ok := ProfileByID(profile)
	if !ok {
		return false
	}
	_, ok = p.bootstrapByName()[name]
	return ok
}

// BootstrapClassDefFor returns profile-scoped class metadata.
func BootstrapClassDefFor(profile, name string) (BootstrapClassInfo, bool) {
	p, ok := ProfileByID(profile)
	if !ok {
		return BootstrapClassInfo{}, false
	}
	c, ok := p.bootstrapByName()[name]
	if !ok {
		return BootstrapClassInfo{}, false
	}
	return classInfoOf(c), true
}

// BootstrapHasMethodFor reports whether the profile bootstrap inventory
// declares the exact member on owner or a superclass.
func BootstrapHasMethodFor(profile, owner, name, descriptor string) bool {
	p, ok := ProfileByID(profile)
	if !ok {
		return false
	}
	return bootstrapHasMethod(p.bootstrapByName(), owner, name, descriptor)
}

// BootstrapHasFieldFor reports whether the profile bootstrap inventory
// declares the exact field on owner or a superclass.
func BootstrapHasFieldFor(profile, owner, name, descriptor string) bool {
	p, ok := ProfileByID(profile)
	if !ok {
		return false
	}
	return bootstrapHasField(p.bootstrapByName(), owner, name, descriptor)
}

// BootstrapMethodsFor lists the profile's intrinsic methods.
func BootstrapMethodsFor(profile string) []BootstrapMethodInfo {
	p, ok := ProfileByID(profile)
	if !ok {
		return nil
	}
	return bootstrapMethodsOf(p.bootstrapInventory())
}

// BootstrapFieldsFor lists the profile's bootstrap fields.
func BootstrapFieldsFor(profile string) []BootstrapFieldInfo {
	p, ok := ProfileByID(profile)
	if !ok {
		return nil
	}
	return bootstrapFieldsOf(p.bootstrapInventory())
}

// BootstrapClassesFor lists the profile's bootstrap classes.
func BootstrapClassesFor(profile string) []BootstrapClassInfo {
	p, ok := ProfileByID(profile)
	if !ok {
		return nil
	}
	return bootstrapClassesOf(p.bootstrapInventory())
}

func (p *Profile) bootstrapInventory() []*bsClass {
	if p.Plugins {
		return gate3BootstrapInventory()
	}
	return bootstrapInventory
}

func (p *Profile) bootstrapByName() map[string]*bsClass {
	if p.Plugins {
		return gate3BootstrapByName()
	}
	return bootstrapByName
}

func classInfoOf(c *bsClass) BootstrapClassInfo {
	return BootstrapClassInfo{Name: c.name, Super: c.super, Interfaces: c.interfaces, Flags: c.flags}
}

func bootstrapHasMethod(byName map[string]*bsClass, owner, name, descriptor string) bool {
	c, ok := byName[owner]
	if !ok {
		return false
	}
	for c != nil {
		for _, m := range c.methods {
			if m.name == name && m.descriptor == descriptor {
				return true
			}
		}
		if c.super == "" {
			break
		}
		c = byName[c.super]
	}
	return false
}

func bootstrapHasField(byName map[string]*bsClass, owner, name, descriptor string) bool {
	c, ok := byName[owner]
	if !ok {
		return false
	}
	for c != nil {
		for _, f := range c.fields {
			if f.name == name && f.descriptor == descriptor {
				return true
			}
		}
		if c.super == "" {
			break
		}
		c = byName[c.super]
	}
	return false
}

func bootstrapMethodsOf(inventory []*bsClass) []BootstrapMethodInfo {
	var out []BootstrapMethodInfo
	for _, c := range inventory {
		for _, m := range c.methods {
			out = append(out, BootstrapMethodInfo{Owner: c.name, Name: m.name, Descriptor: m.descriptor})
		}
	}
	return out
}

func bootstrapFieldsOf(inventory []*bsClass) []BootstrapFieldInfo {
	var out []BootstrapFieldInfo
	for _, c := range inventory {
		for _, f := range c.fields {
			out = append(out, BootstrapFieldInfo{Owner: c.name, Name: f.name, Descriptor: f.descriptor})
		}
	}
	return out
}

func bootstrapClassesOf(inventory []*bsClass) []BootstrapClassInfo {
	out := make([]BootstrapClassInfo, 0, len(inventory))
	for _, c := range inventory {
		out = append(out, classInfoOf(c))
	}
	return out
}
