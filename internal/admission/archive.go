package admission

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"unicode/utf8"

	"catty/internal/classfile"
	"catty/rt"
)

type sourceClass struct {
	class *classfile.Class
	loc   location
}

type archive struct {
	classes  []*sourceClass
	source   map[string][]byte
	deferred int
	digest   string
}

var reservedClassPrefixes = []string{
	"java/", "javax/", "jdk/", "sun/", "com/sun/", "org/w3c/", "org/xml/",
}

func reservedDefinition(name string) bool {
	for _, p := range reservedClassPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return name == "catty/runtime" || strings.HasPrefix(name, "catty/runtime/")
}

func readArchive(data []byte, parent *rt.Loader, parentClasses map[string]*classfile.Class, profile string, limits Limits) (*archive, error) {
	if int64(len(data)) == 0 {
		return nil, failure("ARCHIVE", "", "", "empty JAR")
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, failure("ARCHIVE", "", "", "not a readable JAR: "+err.Error())
	}
	out := &archive{source: map[string][]byte{}}
	if len(zr.File) > limits.Entries {
		return nil, failure("LIMIT", "", "", fmt.Sprintf("entry count %d exceeds %d", len(zr.File), limits.Entries))
	}
	files := map[string]bool{}
	dirs := map[string]bool{}
	classCount := 0
	expanded := int64(0)
	var manifest []byte
	classBodies := map[string][]byte{}
	for _, f := range zr.File {
		name := f.Name
		if err := validateEntryName(name, f.FileInfo().IsDir()); err != nil {
			return nil, failure("ARCHIVE", name, "", err.Error())
		}
		isDir := strings.HasSuffix(name, "/")
		if isDir {
			if files[strings.TrimSuffix(name, "/")] {
				return nil, failure("ARCHIVE", name, "", "file/directory conflict")
			}
			if dirs[name] {
				return nil, failure("ARCHIVE", name, "", "duplicate entry")
			}
			dirs[name] = true
			checkPrefixes(name, files, dirs, true)
			continue
		}
		if dirs[name] || files[name] {
			return nil, failure("ARCHIVE", name, "", "duplicate entry")
		}
		if err := checkPrefixes(name, files, dirs, false); err != nil {
			return nil, failure("ARCHIVE", name, "", err.Error())
		}
		files[name] = true
		if f.Flags&0x1 != 0 {
			return nil, failure("ARCHIVE", name, "", "encrypted entry")
		}
		if f.Method != zip.Store && f.Method != zip.Deflate {
			return nil, failure("ARCHIVE", name, "", fmt.Sprintf("unsupported compression method %d", f.Method))
		}
		mode := f.FileInfo().Mode()
		if mode&os.ModeType != 0 {
			return nil, failure("ARCHIVE", name, "", "disallowed ZIP file type")
		}
		if strings.HasPrefix(name, "META-INF/versions/") {
			return nil, failure("PACKAGE", name, "", "multi-release versioned path")
		}
		if isSignatureEntry(name) {
			return nil, failure("PACKAGE", name, "", "signature metadata")
		}
		if strings.EqualFold(name, "module-info.class") {
			return nil, failure("PACKAGE", name, "", "named module descriptor")
		}
		lower := strings.ToLower(name)
		if strings.HasPrefix(lower, "meta-inf/") && strings.HasSuffix(lower, ".class") {
			return nil, failure("PACKAGE", name, "", "class definition under META-INF")
		}
		if f.UncompressedSize64 > uint64(limits.EntryBytes) {
			return nil, failure("LIMIT", name, "", fmt.Sprintf("entry size %d exceeds %d", f.UncompressedSize64, limits.EntryBytes))
		}
		if int64(f.UncompressedSize64) > limits.ExpandedBytes-expanded {
			return nil, failure("LIMIT", name, "", "expanded data exceeds the JAR ceiling")
		}
		body, err := readEntry(f, limits.EntryBytes)
		if err != nil {
			return nil, failure("ARCHIVE", name, "", err.Error())
		}
		expanded += int64(len(body))
		if expanded > limits.ExpandedBytes {
			return nil, failure("LIMIT", name, "", "expanded data exceeds the JAR ceiling")
		}
		if strings.EqualFold(name, "META-INF/MANIFEST.MF") {
			manifest = body
		}
		if strings.HasSuffix(name, ".class") {
			classCount++
			if classCount > limits.ClassEntries {
				return nil, failure("LIMIT", name, "", fmt.Sprintf("class entry count exceeds %d", limits.ClassEntries))
			}
			classBodies[name] = body
		}
	}
	if classCount == 0 {
		return nil, failure("ARCHIVE", "", "", "JAR contains no class entries")
	}
	if manifest != nil {
		if err := validateManifest(manifest); err != nil {
			return nil, err
		}
	}
	// Sort canonical entry names before class validation so the first class
	// diagnostic is deterministic regardless of central-directory order.
	names := make([]string, 0, len(classBodies))
	for name := range classBodies {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		body := classBodies[name]
		c, err := classfile.Parse(body, name)
		if err != nil {
			code := "CLASS_FORMAT"
			if isProfileParseError(err.Error()) {
				code = "PROFILE"
			}
			return nil, failure(code, name, "", err.Error())
		}
		if c.Name+".class" != name {
			return nil, failure("CLASS_FORMAT", name, c.Name, "entry path does not match this_class")
		}
		if c.Flags&0x8000 != 0 {
			return nil, failure("PROFILE", name, c.Name, "module class")
		}
		if reservedDefinition(c.Name) {
			return nil, failure("PACKAGE", name, c.Name, "reserved class namespace")
		}
		out.classes = append(out.classes, &sourceClass{class: c, loc: location{Entry: name, Class: c.Name}})
		out.source[c.Name] = body
	}
	return out, nil
}

func readEntry(f *zip.File, cap int64) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	lr := io.LimitReader(rc, cap+1)
	body, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("entry read failed: %w", err)
	}
	if int64(len(body)) > cap {
		return nil, fmt.Errorf("entry exceeds %d bytes", cap)
	}
	return body, nil
}

func validateEntryName(name string, isDir bool) error {
	if name == "" {
		return fmt.Errorf("empty entry name")
	}
	if !utf8.ValidString(name) {
		return fmt.Errorf("entry name is not valid UTF-8")
	}
	if strings.ContainsRune(name, 0) {
		return fmt.Errorf("entry name contains NUL")
	}
	if strings.Contains(name, "\\") {
		return fmt.Errorf("entry name contains a backslash")
	}
	if strings.HasPrefix(name, "/") {
		return fmt.Errorf("absolute entry name")
	}
	if len(name) >= 2 && name[1] == ':' && isASCIIAlpha(name[0]) {
		return fmt.Errorf("drive-prefixed entry name")
	}
	trimmed := name
	if strings.HasSuffix(trimmed, "/") {
		if !isDir {
			return fmt.Errorf("directory marker on a non-directory entry")
		}
		trimmed = strings.TrimSuffix(trimmed, "/")
	} else if isDir {
		return fmt.Errorf("directory entry without a trailing slash")
	}
	if trimmed == "" {
		return fmt.Errorf("empty entry path")
	}
	segments := strings.Split(trimmed, "/")
	for _, seg := range segments {
		if seg == "" {
			return fmt.Errorf("repeated path separator")
		}
		if seg == "." || seg == ".." {
			return fmt.Errorf("dot path segment")
		}
	}
	return nil
}

func isASCIIAlpha(b byte) bool { return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') }

func checkPrefixes(name string, files, dirs map[string]bool, isDir bool) error {
	trimmed := strings.TrimSuffix(name, "/")
	parts := strings.Split(trimmed, "/")
	for i := 1; i < len(parts); i++ {
		p := strings.Join(parts[:i], "/")
		if files[p] {
			return fmt.Errorf("directory/file conflict at %s", p)
		}
	}
	if !isDir && dirs[name] {
		return fmt.Errorf("directory/file conflict")
	}
	return nil
}

func isSignatureEntry(name string) bool {
	lower := strings.ToLower(name)
	if !strings.HasPrefix(lower, "meta-inf/") {
		return false
	}
	rest := lower[len("meta-inf/"):]
	base := path.Base(rest)
	switch {
	case strings.HasSuffix(base, ".sf"), strings.HasSuffix(base, ".rsa"),
		strings.HasSuffix(base, ".dsa"), strings.HasSuffix(base, ".ec"):
		return true
	case strings.HasPrefix(base, "sig-"):
		return true
	}
	return false
}

func isProfileParseError(msg string) bool {
	for _, needle := range []string{
		"unsupported class file version", "unsupported opcode", "unsupported constant",
		"unsupported instruction", "invokedynamic", "bootstrap constants",
		"monitor opcode", "legacy subroutine", "reserved opcode",
	} {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}

func sortClasses(classes []*sourceClass) {
	sort.SliceStable(classes, func(i, j int) bool { return classes[i].loc.Entry < classes[j].loc.Entry })
}
