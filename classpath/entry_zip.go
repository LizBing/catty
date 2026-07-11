package classpath

import (
	"archive/zip"
	"io"
)

// zipEntry serves classes from a .jar/.zip file. The central directory is read
// on demand per lookup; for MVP this is only exercised by jar-based classpaths
// (not the MVP test fixtures, which compile to a plain directory).
type zipEntry struct {
	path string
}

func newZipEntry(path string) *zipEntry {
	return &zipEntry{path: path}
}

func (e *zipEntry) ReadClass(name string) ([]byte, Entry, error) {
	r, err := zip.OpenReader(e.path)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	f, err := r.Open(name + ".class")
	if err != nil {
		return nil, nil, errNotFound(name)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	return data, e, nil
}

func (e *zipEntry) String() string { return e.path }
