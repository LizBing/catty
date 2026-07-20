package classpath

import (
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
		// Distinguish true miss from I/O errors (permission, disk, etc.).
		if os.IsNotExist(err) {
			return nil, nil, &ErrNotFound{Name: name}
		}
		return nil, nil, err
	}
	return data, e, nil
}

func (e *dirEntry) String() string { return e.absDir }
