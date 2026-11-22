package admission

import (
	"strconv"
	"strings"

	"catty/rt"
)

type manifestSection struct {
	name  string
	attrs map[string]string
}

// parseManifest implements the JAR manifest's line folding with
// case-insensitive attribute names.
func parseManifest(data []byte) ([]manifestSection, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	var sections []manifestSection
	attrs := map[string]string{}
	sectionName := ""
	flush := func() {
		if len(attrs) == 0 && sectionName == "" {
			return
		}
		sections = append(sections, manifestSection{name: sectionName, attrs: attrs})
		attrs = map[string]string{}
		sectionName = ""
	}
	var pending string
	commit := func(line string) error {
		if line == "" {
			return nil
		}
		colon := strings.IndexByte(line, ':')
		if colon <= 0 {
			return &manifestError{"malformed manifest line"}
		}
		key := strings.ToLower(strings.TrimSpace(line[:colon]))
		value := strings.TrimSpace(line[colon+1:])
		if key == "" {
			return &manifestError{"empty manifest attribute name"}
		}
		attrs[key] = value
		return nil
	}
	for _, raw := range lines {
		if raw == "" {
			if pending != "" {
				if err := commit(pending); err != nil {
					return nil, err
				}
				pending = ""
			}
			flush()
			continue
		}
		if strings.HasPrefix(raw, " ") {
			if pending == "" {
				return nil, &manifestError{"manifest continuation without an attribute"}
			}
			pending += raw[1:]
			continue
		}
		if pending != "" {
			if err := commit(pending); err != nil {
				return nil, err
			}
			pending = ""
		}
		pending = raw
	}
	if pending != "" {
		if err := commit(pending); err != nil {
			return nil, err
		}
	}
	flush()
	return sections, nil
}

type manifestError struct{ msg string }

func (e *manifestError) Error() string { return e.msg }

func validateManifest(data []byte) *rt.PluginError {
	sections, err := parseManifest(data)
	if err != nil {
		return failure("PACKAGE", "META-INF/MANIFEST.MF", "", "malformed manifest: "+err.Error())
	}
	for i, sec := range sections {
		section := "main"
		if i > 0 {
			section = "entry"
		}
		for key, value := range sec.attrs {
			lower := strings.ToLower(value)
			switch key {
			case "multi-release":
				if lower == "true" {
					return locFailure("PACKAGE", location{Entry: "META-INF/MANIFEST.MF", Feature: "multi-release"},
						"manifest declares Multi-Release: true")
				}
			case "class-path":
				if value != "" {
					return locFailure("PACKAGE", location{Entry: "META-INF/MANIFEST.MF", Feature: "class-path"},
						"manifest declares a nonempty Class-Path")
				}
			case "sealed":
				if lower == "true" {
					return locFailure("PACKAGE", location{Entry: "META-INF/MANIFEST.MF", Feature: "sealed"},
						"manifest declares Sealed: true in the "+section+" section")
				}
			}
		}
	}
	return nil
}

var _ = strconv.Itoa
