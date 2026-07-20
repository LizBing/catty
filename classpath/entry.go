package classpath

import (
	"errors"
	"strings"
)

// ErrNotFound is a typed "class not found in this entry" error.
// CompositeEntry treats it as "try the next entry". Real I/O errors,
// corrupt archives, and permission problems are NOT ErrNotFound and
// propagate immediately.
type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	return "catty: class not found: " + e.Name
}

// IsNotFound reports whether err is an ErrNotFound (typed true miss).
func IsNotFound(err error) bool {
	var nf *ErrNotFound
	return errors.As(err, &nf)
}

// Entry is one component of the classpath: a directory or a zip/jar. ReadClass
// returns the bytes of <name>.class (name uses internal slashes) or an error.
// A typed *ErrNotFound means the class does not exist in this entry; any other
// error means the entry itself is broken.
type Entry interface {
	ReadClass(name string) ([]byte, Entry, error)
	String() string
}

// newEntry parses one -cp segment into a dir or zip Entry.
func newEntry(path string) Entry {
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}
	return newDirEntry(path)
}

// CompositeEntry models a multi-segment classpath ("a:b:c").
type CompositeEntry []Entry

func newCompositeEntry(pathList string) CompositeEntry {
	paths := strings.Split(pathList, string(pathListSeparator))
	entries := make([]Entry, 0, len(paths))
	for _, p := range paths {
		entries = append(entries, newEntry(p))
	}
	return entries
}

func (c CompositeEntry) ReadClass(name string) ([]byte, Entry, error) {
	for _, e := range c {
		data, from, err := e.ReadClass(name)
		if err == nil {
			return data, from, nil
		}
		// Typed miss → try the next entry. Real I/O errors propagate immediately.
		if !IsNotFound(err) {
			return nil, nil, err
		}
	}
	return nil, nil, &ErrNotFound{Name: name}
}

func (c CompositeEntry) String() string {
	strs := make([]string, len(c))
	for i, e := range c {
		strs[i] = e.String()
	}
	return strings.Join(strs, string(pathListSeparator))
}
