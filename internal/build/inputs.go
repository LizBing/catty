package build

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"catty/internal/classfile"
	"catty/rt"
)

// Input is one ordered classpath entry with its content digest.
type Input struct {
	Origin string `json:"origin"`
	Digest string `json:"digest"`
}

// Loaded is the deterministic result of reading an ordered classpath.
type Loaded struct {
	Profile *rt.Profile
	Inputs  []Input
	Classes []*classfile.Class
	ByName  map[string]*classfile.Class
	Origins map[string]string
}

var reservedPrefixes = []string{"java/", "javax/", "jdk/", "sun/", "com/sun/", "org/w3c/", "org/xml/"}

func reservedClassName(profile *rt.Profile, name string) bool {
	for _, p := range reservedPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	if profile != nil && profile.Plugins && (name == "catty/runtime" || strings.HasPrefix(name, "catty/runtime/")) {
		return true
	}
	return false
}

// LoadClasspath reads class directories and single-release JARs in order. The
// first definition of a binary name wins and its origin is reported.
func LoadClasspath(entries []string) (*Loaded, error) {
	return LoadClasspathFor(entries, rt.Gate2Profile)
}

// LoadClasspathFor reads an ordered classpath under one profile's admission
// rules. Profile selection therefore affects build-time reserved namespaces as
// well as runtime installation.
func LoadClasspathFor(entries []string, profile *rt.Profile) (*Loaded, error) {
	if profile == nil {
		profile = rt.Gate2Profile
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("catty: empty classpath")
	}
	out := &Loaded{Profile: profile, ByName: map[string]*classfile.Class{}, Origins: map[string]string{}}
	for _, entry := range entries {
		if entry == "" {
			return nil, fmt.Errorf("catty: empty classpath entry")
		}
		abs, err := filepath.Abs(entry)
		if err != nil {
			return nil, fmt.Errorf("catty: classpath entry %s: %w", entry, err)
		}
		info, err := os.Stat(abs)
		if err != nil {
			return nil, fmt.Errorf("catty: classpath entry %s: %w", abs, err)
		}
		digest, err := entryDigest(abs, info)
		if err != nil {
			return nil, err
		}
		out.Inputs = append(out.Inputs, Input{Origin: abs, Digest: digest})
		if info.IsDir() {
			if err := out.loadDirectory(abs); err != nil {
				return nil, err
			}
		} else {
			if err := out.loadJar(abs); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

func entryDigest(path string, info os.FileInfo) (string, error) {
	h := sha256.New()
	if info.IsDir() {
		var files []string
		err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(path, p)
			if err != nil {
				return err
			}
			files = append(files, rel)
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("catty: digest %s: %w", path, err)
		}
		sort.Strings(files)
		for _, rel := range files {
			full := filepath.Join(path, rel)
			data, err := os.ReadFile(full)
			if err != nil {
				return "", fmt.Errorf("catty: digest %s: %w", full, err)
			}
			_, _ = io.WriteString(h, rel)
			_, _ = h.Write([]byte{0})
			_, _ = h.Write(data)
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("catty: digest %s: %w", path, err)
		}
		defer f.Close()
		if _, err := io.Copy(h, f); err != nil {
			return "", fmt.Errorf("catty: digest %s: %w", path, err)
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (l *Loaded) loadDirectory(dir string) error {
	var files []string
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".class") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("catty: reading %s: %w", dir, err)
	}
	sort.Strings(files)
	seen := map[string]bool{}
	for _, path := range files {
		origin, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("catty: reading %s: %w", origin, err)
		}
		if err := l.admit(data, origin, seen); err != nil {
			return err
		}
	}
	return nil
}

func (l *Loaded) loadJar(path string) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("catty: %s is not a readable JAR: %w", path, err)
	}
	defer zr.Close()
	var names []string
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "META-INF/versions/") {
			return fmt.Errorf("catty: %s: multi-release JAR entry %s is unsupported", path, f.Name)
		}
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name, "META-INF/") {
			continue
		}
		if f.Name == "module-info.class" {
			return fmt.Errorf("catty: %s: module-info.class is unsupported", path)
		}
		if strings.HasSuffix(f.Name, ".class") {
			names = append(names, f.Name)
		}
	}
	sort.Strings(names)
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			return fmt.Errorf("catty: %s: duplicate JAR entry %s", path, name)
		}
		seen[name] = true
		reader, err := zr.Open(name)
		if err != nil {
			return fmt.Errorf("catty: %s!/%s: %w", path, name, err)
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			return fmt.Errorf("catty: %s!/%s: %w", path, name, err)
		}
		origin := path + "!/" + name
		if err := l.admit(data, origin, nil); err != nil {
			return err
		}
	}
	return nil
}

func (l *Loaded) admit(data []byte, origin string, sameEntry map[string]bool) error {
	c, err := classfile.Parse(data, origin)
	if err != nil {
		return fmt.Errorf("%s: %w", origin, err)
	}
	if reservedClassName(l.Profile, c.Name) {
		return fmt.Errorf("%s: class %s is in the reserved bootstrap namespace", origin, c.Name)
	}
	if _, exists := l.ByName[c.Name]; exists {
		// First classpath entry wins; a later definition is shadowed.
		return nil
	}
	if sameEntry != nil {
		if sameEntry[c.Name] {
			return fmt.Errorf("%s: duplicate definition of %s in one classpath entry", origin, c.Name)
		}
		sameEntry[c.Name] = true
	}
	l.ByName[c.Name] = c
	l.Origins[c.Name] = origin
	l.Classes = append(l.Classes, c)
	return nil
}
