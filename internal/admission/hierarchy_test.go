package admission

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"catty/rt"
)

func u2bytes(v uint16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], v)
	return b[:]
}

func u4bytes(v uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], v)
	return b[:]
}

type cpBuilder struct {
	entries []byte
	next    uint16
	classes map[string]uint16
}

func newCP() *cpBuilder { return &cpBuilder{next: 1, classes: map[string]uint16{}} }

func (c *cpBuilder) utf8(s string) uint16 {
	idx := c.next
	c.next++
	c.entries = append(c.entries, 1)
	c.entries = append(c.entries, u2bytes(uint16(len(s)))...)
	c.entries = append(c.entries, []byte(s)...)
	return idx
}

func (c *cpBuilder) class(name string) uint16 {
	if idx, ok := c.classes[name]; ok {
		return idx
	}
	nameIdx := c.utf8(name)
	idx := c.next
	c.next++
	c.entries = append(c.entries, 7)
	c.entries = append(c.entries, u2bytes(nameIdx)...)
	c.classes[name] = idx
	return idx
}

// minimalClass builds a class-file-version-65 class with no fields, methods or
// attributes. It is enough to exercise the structural, hierarchy and profile
// checks at the admission boundary.
func minimalClass(name, super string, interfaces []string, flags uint16) []byte {
	cp := newCP()
	thisIdx := cp.class(name)
	superIdx := cp.class(super)
	ifaceIdxs := make([]uint16, 0, len(interfaces))
	for _, in := range interfaces {
		ifaceIdxs = append(ifaceIdxs, cp.class(in))
	}
	var out bytes.Buffer
	out.Write([]byte{0xCA, 0xFE, 0xBA, 0xBE})
	out.Write(u2bytes(0))
	out.Write(u2bytes(65))
	out.Write(u2bytes(cp.next))
	out.Write(cp.entries)
	out.Write(u2bytes(flags))
	out.Write(u2bytes(thisIdx))
	out.Write(u2bytes(superIdx))
	out.Write(u2bytes(uint16(len(ifaceIdxs))))
	for _, idx := range ifaceIdxs {
		out.Write(u2bytes(idx))
	}
	out.Write(u2bytes(0)) // fields
	out.Write(u2bytes(0)) // methods
	out.Write(u2bytes(0)) // attributes
	return out.Bytes()
}

func depthJar(prefix string, classes int) []byte {
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	// Ancestor-first archive order is the reviewer reproducer: C000 appears
	// before every descendant and is completed long before C128 is visited.
	for i := 0; i < classes; i++ {
		name := fmt.Sprintf("%s/C%03d", prefix, i)
		super := "java/lang/Object"
		if i > 0 {
			super = fmt.Sprintf("%s/C%03d", prefix, i-1)
		}
		data := minimalClass(name, super, nil, rt.AccPublic)
		w, err := zw.Create(name + ".class")
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}
	return out.Bytes()
}

func TestAdmissionHierarchyDepthExactBoundary(t *testing.T) {
	parent := parentLoader(t)
	limits := DefaultLimits()
	res, err := AdmitWithLimits(depthJar("chain", 128), parent, rt.Gate3ProfileID, limits)
	if err != nil {
		t.Fatalf("128-edge chain rejected: %v", err)
	}
	if len(res.Defined) != 128 {
		t.Fatalf("128-edge chain defined %d classes", len(res.Defined))
	}
	_, err = AdmitWithLimits(depthJar("chain", 129), parent, rt.Gate3ProfileID, limits)
	if err == nil {
		t.Fatal("129-edge chain admitted")
	}
	pe, ok := err.(*rt.PluginError)
	if !ok || pe.Code != "LIMIT" || pe.Limit != "hierarchyDepth" {
		t.Fatalf("129-edge diagnostic: %#v", err)
	}
}

func parentChainLoader(t *testing.T, n int) *rt.Loader {
	t.Helper()
	r := rt.NewRuntime(&rt.Image{ABI: rt.Gate3ABI, Profile: rt.Gate3ProfileID, Plugins: true})
	if err := r.Install(); err != nil {
		t.Fatal(err)
	}
	object, _ := r.BootstrapLoader.Lookup("java/lang/Object")
	loader := &rt.Loader{ID: 2, Name: "span-parent", Runtime: r, Parent: r.BootstrapLoader,
		Classes: map[string]*rt.Class{}}
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("span/C%03d", i)
		super := object
		if i > 0 {
			super = loader.Classes[fmt.Sprintf("span/C%03d", i-1)]
		}
		loader.Classes[name] = &rt.Class{Name: name, Super: super, Loader: loader,
			Flags: rt.AccPublic, MethodIndex: map[string]*rt.Method{}}
	}
	return loader
}

func TestAdmissionHierarchyDepthSpansParentLoader(t *testing.T) {
	parent := parentChainLoader(t, 120)
	limits := DefaultLimits()
	child := func(count int) []byte {
		var out bytes.Buffer
		zw := zip.NewWriter(&out)
		for n := 120; n < 120+count; n++ {
			name := fmt.Sprintf("span/C%03d", n)
			super := fmt.Sprintf("span/C%03d", n-1)
			w, _ := zw.Create(name + ".class")
			_, _ = w.Write(minimalClass(name, super, nil, rt.AccPublic))
		}
		_ = zw.Close()
		return out.Bytes()
	}
	if _, err := AdmitWithLimits(child(8), parent, rt.Gate3ProfileID, limits); err != nil {
		t.Fatalf("128-edge parent-spanning chain rejected: %v", err)
	}
	_, err := AdmitWithLimits(child(9), parent, rt.Gate3ProfileID, limits)
	pe, ok := err.(*rt.PluginError)
	if !ok || pe.Code != "LIMIT" || pe.Limit != "hierarchyDepth" {
		t.Fatalf("129-edge parent-spanning chain: %#v", err)
	}
}

func TestAdmissionHierarchyCycleRejected(t *testing.T) {
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	for _, pair := range [][2]string{{"cycle/A", "cycle/B"}, {"cycle/B", "cycle/A"}} {
		w, _ := zw.Create(pair[0] + ".class")
		_, _ = w.Write(minimalClass(pair[0], pair[1], nil, rt.AccPublic))
	}
	_ = zw.Close()
	_, err := AdmitWithLimits(out.Bytes(), parentLoader(t), rt.Gate3ProfileID, DefaultLimits())
	pe, ok := err.(*rt.PluginError)
	if !ok || pe.Code != "VERIFY" || pe.Class == "" {
		t.Fatalf("cycle diagnostic: %#v", err)
	}
}
