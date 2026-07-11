package classpath

import (
	"errors"
	"os"
	"path/filepath"
)

// dirEntry serves a classpath directory. Class java.lang.Object maps to
// <dir>/java/lang/Object.class.
type dirEntry struct {
	absDir string
}

func newDirEntry(path string) *dirEntry {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return &dirEntry{absDir: abs}
}

func (e *dirEntry) ReadClass(name string) ([]byte, Entry, error) {
	full := filepath.Join(e.absDir, filepath.FromSlash(name+".class"))
	data, err := os.ReadFile(full)
	if err != nil {
		return nil, nil, err
	}
	return data, e, nil
}

func (e *dirEntry) String() string { return e.absDir }

// errNotFound is the canonical miss; callers (CompositeEntry) treat any error
// as "try the next entry", so a distinct error just aids diagnostics.
func errNotFound(name string) error {
	return errors.New("catty: class not found: " + name)
}
