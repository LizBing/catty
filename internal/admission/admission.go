// Package admission implements Gate 3 bounded JAR validation and staging. It
// parses and verifies a candidate with the shared class-file reader and typed
// verifier, then returns immutable definitions for the runtime to commit in a
// single transaction. It never mutates runtime state and never extracts or
// fetches dependencies.
package admission

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"catty/internal/classfile"
	"catty/internal/verify"
	"catty/rt"
)

// Limits are one admission transaction's ceilings. Tests may lower them; the
// production values come from the published profile.
type Limits struct {
	CompressedBytes      int64
	ExpandedBytes        int64
	EntryBytes           int64
	Entries              int
	ClassEntries         int
	HierarchyDepth       int
	VerificationWork     int64
	Loaders              int
	CumulativeClassBytes int64
}

// DefaultLimits returns the Gate 3 profile ceilings.
func DefaultLimits() Limits {
	l := rt.Gate3PluginLimits
	return Limits{
		CompressedBytes:      l.CompressedBytes,
		ExpandedBytes:        l.ExpandedBytes,
		EntryBytes:           l.EntryBytes,
		Entries:              l.Entries,
		ClassEntries:         l.ClassEntries,
		HierarchyDepth:       l.HierarchyDepth,
		VerificationWork:     l.VerificationWork,
		Loaders:              l.Loaders,
		CumulativeClassBytes: l.CumulativeClassBytes,
	}
}

func failure(code, entry, class, message string) *rt.PluginError {
	return &rt.PluginError{Code: code, Entry: entry, Class: class, Message: message}
}

func locFailure(code string, loc location, message string) *rt.PluginError {
	err := &rt.PluginError{
		Code: code, Entry: loc.Entry, Class: loc.Class, Method: loc.Method,
		Descriptor: loc.Descriptor, PC: loc.PC, Feature: loc.Feature, Limit: loc.Limit,
		Message: message,
	}
	return err
}

type location struct {
	Entry      string
	Class      string
	Method     string
	Descriptor string
	PC         int
	Feature    string
	Limit      string
}

// Admit validates one candidate under the Gate 3 profile.
func Admit(data []byte, parent *rt.Loader, profile string) (*rt.AdmissionResult, error) {
	return AdmitWithLimits(data, parent, profile, DefaultLimits())
}

// AdmitWithLimits is Admit with explicit ceilings, used by the Go tests.
func AdmitWithLimits(data []byte, parent *rt.Loader, profile string, limits Limits) (*rt.AdmissionResult, error) {
	if parent == nil {
		return nil, failure("VERIFY", "", "", "nil parent loader")
	}
	if profile != rt.Gate3ProfileID {
		return nil, failure("PROFILE", "", "", "plugin admission requires the Gate 3 profile")
	}
	if limits.CompressedBytes > 0 && int64(len(data)) > limits.CompressedBytes {
		return nil, failure("LIMIT", "", "", fmt.Sprintf("compressed input exceeds %d bytes", limits.CompressedBytes))
	}
	if limits.ExpandedBytes <= 0 || limits.EntryBytes <= 0 || limits.Entries <= 0 ||
		limits.ClassEntries <= 0 || limits.HierarchyDepth <= 0 || limits.VerificationWork <= 0 {
		return nil, failure("LIMIT", "", "", "invalid admission limits")
	}
	sum := sha256.Sum256(data)
	digest := hex.EncodeToString(sum[:])

	parentClasses := parentClassTable(parent)
	archive, err := readArchive(data, parent, parentClasses, profile, limits)
	if err != nil {
		return nil, err
	}
	archive.digest = digest
	candidate := archive.classes
	if len(candidate) == 0 {
		return nil, failure("ARCHIVE", "", "", "JAR contains no class entries")
	}

	// Sort canonical entry names for deterministic first diagnostics.
	sortClasses(candidate)

	byName := map[string]*classfile.Class{}
	var effective []*classfile.Class
	var shadowed []string
	for _, sc := range candidate {
		c := sc.class
		if _, ok := byName[c.Name]; ok {
			return nil, locFailure("CLASS_FORMAT", sc.loc, "duplicate internal class name "+c.Name)
		}
		byName[c.Name] = c
		if _, ok := parentClasses[c.Name]; ok {
			shadowed = append(shadowed, c.Name)
			continue
		}
		effective = append(effective, c)
	}
	if len(effective) == 0 {
		// A JAR of entirely shadowed classes still publishes a loader.
	}

	ctxClasses := map[string]*classfile.Class{}
	for name, c := range parentClasses {
		ctxClasses[name] = c
	}
	for name, c := range byName {
		if _, ok := ctxClasses[name]; !ok {
			ctxClasses[name] = c
		}
	}
	work := &verify.Work{Limit: limits.VerificationWork}
	ctx := &verify.Context{Classes: ctxClasses, Profile: profile, Work: work}
	ctx.MemberFlags = func(owner, name, descriptor string) (uint16, string, bool) {
		return resolveMemberFlags(owner, name, descriptor, byName, parentClasses)
	}
	ctx.RuntimePackage = func(a, b string) bool {
		_, aCandidate := byName[a]
		_, bCandidate := byName[b]
		if aCandidate != bCandidate {
			// Candidate and parent classes are defined by different loaders.
			return false
		}
		return packageOf(a) == packageOf(b)
	}

	// Structural/profile checks and hierarchy resolution apply to every parsed
	// entry, including definitions the parent shadows.
	if err := checkCandidate(archive, candidate, parentClasses, byName, profile, limits, work); err != nil {
		return nil, err
	}
	if err := checkHierarchy(effective, byName, parentClasses, limits, work); err != nil {
		return nil, err
	}
	if err := checkNests(effective, byName, parentClasses, parent, limits.HierarchyDepth); err != nil {
		return nil, err
	}

	// The shared typed verifier consumes the effective set. Shadowed entries
	// receive structural checks above but are not definitions.
	for _, c := range effective {
		if err := verify.VerifyClass(c, ctx); err != nil {
			return nil, classifyVerifyError(c, err)
		}
	}

	if err := checkBootstrapReferences(effective, byName, parentClasses, profile); err != nil {
		return nil, err
	}

	defs := make([]*rt.ClassDef, 0, len(effective))
	interpreted := 0
	classBytes := int64(0)
	for _, c := range effective {
		def, err := define(c)
		if err != nil {
			return nil, locFailure("CLASS_FORMAT", location{Class: c.Name}, err.Error())
		}
		for _, m := range def.Methods {
			if m.Engine == "interp" {
				interpreted++
			}
		}
		defs = append(defs, def)
		if src, ok := archive.source[c.Name]; ok {
			classBytes += int64(len(src))
		}
	}
	if limits.CumulativeClassBytes > 0 && classBytes > limits.CumulativeClassBytes {
		return nil, failure("LIMIT", "", "", "candidate class bytes exceed the runtime cumulative limit")
	}
	defined := make([]string, 0, len(effective))
	for _, c := range effective {
		defined = append(defined, c.Name)
	}
	var deferredRefs []rt.DeferredReference
	for _, c := range effective {
		for _, m := range c.Methods {
			pcOf := map[int32]int32{}
			for _, insn := range m.Insns {
				if _, ok := pcOf[insn.Ref]; !ok {
					pcOf[insn.Ref] = insn.PC
				}
			}
			for index, ref := range m.Refs {
				kind := ""
				owner, name, descriptor := ref.Owner, ref.Name, ref.Descriptor
				switch ref.Kind {
				case rt.RefField:
					kind = "field"
				case rt.RefMethod:
					kind = "method"
				case rt.RefClass:
					kind = "class"
					owner, name, descriptor = ref.Class, "", ""
				default:
					continue
				}
				visible := true
				if kind == "class" {
					_, visible = lookupMeta(ref.Class, byName, parentClasses)
				} else {
					_, visible = lookupMeta(ref.Owner, byName, parentClasses)
				}
				if visible || owner == "" {
					continue
				}
				deferredRefs = append(deferredRefs, rt.DeferredReference{
					Kind: kind, Owner: owner, Name: name, Descriptor: descriptor,
					Referrer: c.Name, ReferrerMethod: m.Name, ReferrerDescriptor: m.Descriptor,
					PC: int(pcOf[int32(index)]),
				})
			}
		}
	}
	return &rt.AdmissionResult{
		DeferredRefs: deferredRefs,
		Classes:      defs, Defined: defined, Shadowed: shadowed,
		Interpreted: interpreted, Deferred: len(deferredRefs),
		Digest: digest, ClassBytes: classBytes, WorkUsed: work.Used,
	}, nil
}

func classifyVerifyError(c *classfile.Class, err error) *rt.PluginError {
	msg := err.Error()
	code := "VERIFY"
	var workErr *verify.WorkError
	if errors.As(err, &workErr) {
		code = "LIMIT"
	} else if isDependencyMessage(msg) {
		code = "DEPENDENCY"
	} else if isProfileParseError(msg) {
		code = "PROFILE"
	}
	loc := location{Class: c.Name, Limit: ""}
	if code == "LIMIT" {
		loc.Limit = "verificationWork"
	}
	perr := locFailure(code, loc, msg)
	perr.Method, perr.Descriptor, perr.PC = verifyLocation(c.Name, msg)
	if workErr != nil {
		perr.Feature = workErr.Operation
	}
	return perr
}

// verifyLocation recovers the method and bytecode offset from the verifier's
// stable "Class.method(desc): pc N: ..." diagnostic prefix.
func verifyLocation(className, msg string) (string, string, int) {
	prefix := className + "."
	if !strings.HasPrefix(msg, prefix) {
		return "", "", 0
	}
	rest := strings.TrimPrefix(msg, prefix)
	colon := strings.Index(rest, ":")
	if colon < 0 {
		return "", "", 0
	}
	signature := rest[:colon]
	open := strings.IndexByte(signature, '(')
	if open < 0 {
		return "", "", 0
	}
	method, descriptor := signature[:open], signature[open:]
	pc := 0
	if index := strings.Index(msg, "pc "); index >= 0 {
		start := index + 3
		end := start
		for end < len(msg) && msg[end] >= '0' && msg[end] <= '9' {
			end++
		}
		if end > start {
			pc, _ = strconv.Atoi(msg[start:end])
		}
	}
	return method, descriptor, pc
}

func isDependencyMessage(msg string) bool {
	for _, needle := range []string{
		"not available", "unknown superclass", "unknown interface", "unknown class",
		"unknown component",
		"catch type", "cannot resolve", "missing",
	} {
		if containsFold(msg, needle) {
			return true
		}
	}
	return false
}

func containsFold(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexFold(s, sub) >= 0)
}

func indexFold(s, sub string) int {
	// Small ASCII-insensitive index; error text is diagnostic only.
	for i := 0; i+len(sub) <= len(s); i++ {
		ok := true
		for j := 0; j < len(sub); j++ {
			a, b := s[i+j], sub[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}
