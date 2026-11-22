package rt

import (
	"encoding/json"
	"fmt"
)

// DecodeImage parses an embedded image specification and attaches the
// executable's compiled entry table. Generated applications use this at
// process startup; the returned image is still subject to Install's ABI and
// profile checks.
func DecodeImage(data []byte, entries []Entry) (*Image, error) {
	var img Image
	if err := json.Unmarshal(data, &img); err != nil {
		return nil, fmt.Errorf("catty: cannot decode embedded image: %w", err)
	}
	if err := img.AttachEntries(entries); err != nil {
		return nil, err
	}
	return &img, nil
}

// AttachEntries installs compiled function values by image table index. It is
// used by DecodeImage and by tests that construct images directly.
func (img *Image) AttachEntries(entries []Entry) error {
	for _, c := range img.Classes {
		if c == nil {
			return fmt.Errorf("catty: image contains a nil class definition")
		}
		for _, m := range c.Methods {
			if m == nil {
				return fmt.Errorf("catty: class %s contains a nil method definition", c.Name)
			}
			if m.EntryIndex < 0 {
				continue
			}
			if m.EntryIndex >= len(entries) {
				return fmt.Errorf("catty: method %s.%s%s refers to entry %d of %d",
					m.Owner, m.Name, m.Descriptor, m.EntryIndex, len(entries))
			}
			if entries[m.EntryIndex] == nil {
				return fmt.Errorf("catty: method %s.%s%s has a nil compiled entry",
					m.Owner, m.Name, m.Descriptor)
			}
			m.Entry = entries[m.EntryIndex]
		}
	}
	return nil
}
