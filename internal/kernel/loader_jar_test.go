package kernel

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestJarClasspathLoading pins DEBT-0008 closure: a .jar classpath entry
// is indexed at loader construction and its classes load (and verify)
// exactly like directory entries. The jar is built in-test from the real
// HelloWorld fixture bytes.
func TestJarClasspathLoading(t *testing.T) {
	fixture, err := os.ReadFile(filepath.Join("..", "..", "testdata", "HelloWorld.class"))
	if err != nil {
		t.Fatal(err)
	}

	jarPath := "/tmp/hello-test.jar"
	_ = os.Remove(jarPath)
	if err := writeZip(jarPath, map[string][]byte{
		"HelloWorld.class": fixture,
	}); err != nil {
		t.Fatal(err)
	}

	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	l := NewClassPathLoader(k, []string{jarPath})

	c, err := l.Load("HelloWorld")
	if err != nil {
		t.Fatalf("load from jar: %v", err)
	}
	if c == nil || c.Name != "HelloWorld" {
		t.Fatalf("unexpected class %v", c)
	}
	// Second load resolves from the registry.
	if c2, err := l.Load("HelloWorld"); err != nil || c2 != c {
		t.Fatalf("registry re-load mismatch: %v %v", c2, err)
	}
}

// TestJarMissingEntrySkipped pins that an unreadable/missing .jar entry is
// skipped silently and lookup falls through to later entries.
func TestJarMissingEntrySkipped(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	dir := t.TempDir()
	fixture, _ := os.ReadFile(filepath.Join("..", "..", "testdata", "HelloWorld.class"))
	if err := os.WriteFile(filepath.Join(dir, "HelloWorld.class"), fixture, 0o644); err != nil {
		t.Fatal(err)
	}
	l := NewClassPathLoader(k, []string{
		filepath.Join(t.TempDir(), "no-such.jar"), // missing jar skipped
		dir,                                       // directory fallback wins
	})
	c, err := l.Load("HelloWorld")
	if err != nil {
		t.Fatalf("fallback load failed: %v", err)
	}
	if !strings.HasSuffix(c.Name, "HelloWorld") {
		t.Fatalf("wrong class %s", c.Name)
	}
}

func writeZip(path string, files map[string][]byte) error {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, data := range files {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
