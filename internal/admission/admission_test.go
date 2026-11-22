package admission

import (
	"archive/zip"
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"testing"

	"catty/rt"
)

//go:embed testdata/*.jar
var fixtures embed.FS

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := fixtures.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func parentLoader(t *testing.T) *rt.Loader {
	t.Helper()
	r := rt.NewRuntime(&rt.Image{ABI: rt.Gate3ABI, Profile: rt.Gate3ProfileID, Plugins: true})
	if err := r.Install(); err != nil {
		t.Fatal(err)
	}
	return r.BootstrapLoader
}

func acceptable(t *testing.T, data []byte) *rt.AdmissionResult {
	t.Helper()
	res, err := Admit(data, parentLoader(t), rt.Gate3ProfileID)
	if err != nil {
		t.Fatalf("admission failed: %v", err)
	}
	if res == nil {
		t.Fatal("nil result")
	}
	return res
}

func expectCode(t *testing.T, data []byte, want string) {
	t.Helper()
	_, err := Admit(data, parentLoader(t), rt.Gate3ProfileID)
	if err == nil {
		t.Fatalf("expected %s, got acceptance", want)
	}
	pe, ok := err.(*rt.PluginError)
	if !ok {
		t.Fatalf("expected *rt.PluginError, got %T: %v", err, err)
	}
	if pe.Code != want {
		t.Fatalf("expected code %s, got %s (%v)", want, pe.Code, err)
	}
}

func TestAdmitValidStoredAndDeflated(t *testing.T) {
	for _, name := range []string{"valid.jar", "valid-stored.jar", "valid-deflated.jar"} {
		res := acceptable(t, fixture(t, name))
		if len(res.Defined) == 0 || res.Interpreted == 0 {
			t.Fatalf("%s: definitions=%v interpreted=%d", name, res.Defined, res.Interpreted)
		}
		if res.ClassBytes <= 0 {
			t.Fatalf("%s: class bytes not counted", name)
		}
	}
}

func TestAdmitRejectsArchiveTransport(t *testing.T) {
	expectCode(t, []byte("not a zip"), "ARCHIVE")
	expectCode(t, nil, "ARCHIVE")
	valid := fixture(t, "valid.jar")
	expectCode(t, valid[:len(valid)-20], "ARCHIVE")
}

func TestAdmitReservedNamespace(t *testing.T) {
	plain := fixture(t, "valid-stored.jar")
	zr, err := zip.NewReader(bytes.NewReader(plain), int64(len(plain)))
	if err != nil {
		t.Fatal(err)
	}
	var classBytes []byte
	for _, f := range zr.File {
		rc, _ := f.Open()
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(rc)
		_ = rc.Close()
		classBytes = buf.Bytes()
	}
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	w, _ := zw.Create("catty/runtime/Evil.class")
	_, _ = w.Write(classBytes)
	_ = zw.Close()
	// The class's internal name is ext/Plain, so the entry/name mismatch would
	// also reject. Either way, the candidate must not publish.
	expectCode(t, out.Bytes(), "CLASS_FORMAT")
}

func TestAdmitLimits(t *testing.T) {
	plain := fixture(t, "valid-stored.jar")
	zr, _ := zip.NewReader(bytes.NewReader(plain), int64(len(plain)))
	var classBytes []byte
	for _, f := range zr.File {
		rc, _ := f.Open()
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(rc)
		_ = rc.Close()
		classBytes = buf.Bytes()
	}
	build := func(entries map[string][]byte) []byte {
		var out bytes.Buffer
		zw := zip.NewWriter(&out)
		names := make([]string, 0, len(entries))
		for name := range entries {
			names = append(names, name)
		}
		// Determinism is not required by the test; map order only changes the
		// first diagnostic among equally invalid archives.
		for _, name := range names {
			w, _ := zw.Create(name)
			_, _ = w.Write(entries[name])
		}
		_ = zw.Close()
		return out.Bytes()
	}
	// Entry ceiling: exactly 4 MiB is accepted, one byte more is LIMIT.
	exact := build(map[string][]byte{"ext/Plain.class": classBytes, "pad.bin": make([]byte, 4<<20)})
	acceptable(t, exact)
	over := build(map[string][]byte{"ext/Plain.class": classBytes, "pad.bin": make([]byte, (4<<20)+1)})
	expectCode(t, over, "LIMIT")

	limits := DefaultLimits()
	limits.VerificationWork = 1
	if _, err := AdmitWithLimits(plain, parentLoader(t), rt.Gate3ProfileID, limits); err == nil {
		t.Fatal("verification work ceiling was not enforced")
	} else if pe, ok := err.(*rt.PluginError); !ok || pe.Code != "LIMIT" {
		t.Fatalf("expected LIMIT, got %v", err)
	}
}

func FuzzStageJar(f *testing.F) {
	f.Add(fixtureFile("valid-stored.jar"))
	f.Add(fixtureFile("valid-deflated.jar"))
	f.Add([]byte("not a zip"))
	f.Add([]byte{0x50, 0x4b, 0x03, 0x04})
	parent := &rt.Loader{ID: 0, Name: "bootstrap", Classes: map[string]*rt.Class{}}
	// A real bootstrap loader is needed for the synthetic parent chain.
	r := rt.NewRuntime(&rt.Image{ABI: rt.Gate3ABI, Profile: rt.Gate3ProfileID, Plugins: true})
	if err := r.Install(); err != nil {
		f.Fatal(err)
	}
	parent = r.BootstrapLoader
	limits := DefaultLimits()
	// Keep fuzz inputs well inside the archive reader's production ceilings so
	// the target spends its time in validation, not in a multi-megabyte ZIP
	// directory allocation.
	limits.CompressedBytes = 1 << 20
	limits.ExpandedBytes = 4 << 20
	limits.EntryBytes = 1 << 20
	limits.Entries = 128
	limits.ClassEntries = 32
	limits.VerificationWork = 1 << 20
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			return
		}
		res, err := AdmitWithLimits(data, parent, rt.Gate3ProfileID, limits)
		if err == nil && res == nil {
			t.Fatal("success without definitions")
		}
		if err != nil {
			if _, ok := err.(*rt.PluginError); !ok {
				t.Fatalf("untyped admission error %T: %v", err, err)
			}
		}
	})
}

func fixtureFile(name string) []byte {
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		panic(err)
	}
	return data
}
