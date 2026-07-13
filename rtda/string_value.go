package rtda

// StringValue is the canonical, immutable kernel backing of every Java String
// value. It holds a sequence of UTF-16 code units and implements all of the
// Java 25 String semantics required by the current supported surface.
//
// Every constructor and accessor that exposes units to the caller performs a
// defensive copy — no mutable alias to the backing can escape. Substring and
// concat results are independent values.
//
// This type is the Go-native representation decided by ADR-0027. It is placed
// in Object.Extra() for the current synthetic java/lang/String facade; a future
// ordinary java.base String implementation may store it differently.
type StringValue struct {
	units []uint16
}

// NewStringValue creates a StringValue from a UTF-16 code-unit sequence. The
// slice is defensively copied so the caller retains no alias — no mutable
// backing can escape.
func NewStringValue(units []uint16) *StringValue {
	cp := make([]uint16, len(units))
	copy(cp, units)
	return &StringValue{units: cp}
}

// Len returns the number of UTF-16 code units — Java's String.length().
func (sv *StringValue) Len() int { return len(sv.units) }

// CharAt returns the UTF-16 code unit at index i. Caller must bounds-check.
func (sv *StringValue) CharAt(i int) uint16 { return sv.units[i] }

// IsEmpty reports whether the string has zero code units.
func (sv *StringValue) IsEmpty() bool { return len(sv.units) == 0 }

// Units returns a defensive copy of the code-unit sequence. Callers that only
// need to read (not retain) the units should iterate in place without calling
// this.
func (sv *StringValue) Units() []uint16 {
	cp := make([]uint16, len(sv.units))
	copy(cp, sv.units)
	return cp
}

// HashCode computes Java's String.hashCode(): s[0]*31^(n-1) + ... + s[n-1]
// over UTF-16 code units, with s[i] treated as an unsigned 16-bit value.
func (sv *StringValue) HashCode() int32 {
	h := int32(0)
	for _, u := range sv.units {
		h = h*31 + int32(u)
	}
	return h
}

// Equals implements Java String.equals: two strings are equal iff they have
// the same length and identical code units at every position.
func (sv *StringValue) Equals(other *StringValue) bool {
	if len(sv.units) != len(other.units) {
		return false
	}
	for i, u := range sv.units {
		if u != other.units[i] {
			return false
		}
	}
	return true
}

// CompareTo implements Java String.compareTo lexicographically over unsigned
// 16-bit code units.
func (sv *StringValue) CompareTo(other *StringValue) int {
	n := len(sv.units)
	if len(other.units) < n {
		n = len(other.units)
	}
	for i := 0; i < n; i++ {
		a, b := sv.units[i], other.units[i]
		if a != b {
			return int(a) - int(b)
		}
	}
	return len(sv.units) - len(other.units)
}

// Substring returns a new StringValue for units[begin:end]. begin is
// inclusive, end is exclusive. Caller must validate bounds.
func (sv *StringValue) Substring(begin, end int) *StringValue {
	cp := make([]uint16, end-begin)
	copy(cp, sv.units[begin:end])
	return &StringValue{units: cp}
}

// Concat returns a new StringValue that is the concatenation of sv and other.
func (sv *StringValue) Concat(other *StringValue) *StringValue {
	cp := make([]uint16, len(sv.units)+len(other.units))
	copy(cp, sv.units)
	copy(cp[len(sv.units):], other.units)
	return &StringValue{units: cp}
}

// IndexOf implements Java String.indexOf(int) with code-point rules: if ch is
// a supplementary code point (>= 0x10000), search for its high-then-low
// surrogate pair; otherwise search for the single code unit.
func (sv *StringValue) IndexOf(ch int) int {
	if ch < 0 || ch > 0x10FFFF {
		return -1
	}
	if ch >= 0x10000 {
		// Supplementary: search for surrogate pair.
		hi := uint16((ch-0x10000)>>10) + 0xD800
		lo := uint16((ch-0x10000)&0x3FF) + 0xDC00
		for i := 0; i < len(sv.units)-1; i++ {
			if sv.units[i] == hi && sv.units[i+1] == lo {
				return i
			}
		}
		return -1
	}
	// BMP: search for the single code unit.
	for i, u := range sv.units {
		if u == uint16(ch) {
			return i
		}
	}
	return -1
}

// StartsWith reports whether sv starts with prefix.
func (sv *StringValue) StartsWith(prefix *StringValue) bool {
	if len(prefix.units) > len(sv.units) {
		return false
	}
	for i, u := range prefix.units {
		if sv.units[i] != u {
			return false
		}
	}
	return true
}

// EndsWith reports whether sv ends with suffix.
func (sv *StringValue) EndsWith(suffix *StringValue) bool {
	if len(suffix.units) > len(sv.units) {
		return false
	}
	off := len(sv.units) - len(suffix.units)
	for i, u := range suffix.units {
		if sv.units[off+i] != u {
			return false
		}
	}
	return true
}

// ToCharArray returns a fresh copy of the code units as a new []uint16.
func (sv *StringValue) ToCharArray() []uint16 {
	cp := make([]uint16, len(sv.units))
	copy(cp, sv.units)
	return cp
}

// GoString converts the UTF-16 code-unit sequence to a Go string suitable for
// host output (stdout, stderr, files). Valid surrogate pairs are decoded into
// their UTF-8 scalar representation; each unpaired surrogate (isolated high
// or low) is encoded as '?' (0x3f), matching the pinned Temurin 25.0.3 UTF-8
// PrintStream observation.
//
// This is an output-boundary adapter, never a second canonical Java String
// representation. Java-to-Java paths must use the code units directly.
func (sv *StringValue) GoString() string {
	if len(sv.units) == 0 {
		return ""
	}
	// Fast path: pure ASCII needs no conversion.
	ascii := true
	for _, u := range sv.units {
		if u >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		b := make([]byte, len(sv.units))
		for i, u := range sv.units {
			b[i] = byte(u)
		}
		return string(b)
	}

	// Slow path: encode valid surrogate pairs as UTF-8, lone surrogates as '?'.
	var out []byte
	for i := 0; i < len(sv.units); i++ {
		u := sv.units[i]
		if u >= 0xD800 && u <= 0xDBFF {
			// High surrogate: look for a following low surrogate.
			if i+1 < len(sv.units) && sv.units[i+1] >= 0xDC00 && sv.units[i+1] <= 0xDFFF {
				hi, lo := u, sv.units[i+1]
				cp := 0x10000 + (int(hi)-0xD800)<<10 | (int(lo) - 0xDC00)
				out = appendRune(out, rune(cp))
				i++ // consumed the low surrogate
				continue
			}
			// Lone high surrogate → '?'.
			out = append(out, '?')
			continue
		}
		if u >= 0xDC00 && u <= 0xDFFF {
			// Lone low surrogate → '?'.
			out = append(out, '?')
			continue
		}
		// BMP character — encode as UTF-8.
		out = appendRune(out, rune(u))
	}
	return string(out)
}

// appendRune encodes r as UTF-8 and appends to out.
func appendRune(out []byte, r rune) []byte {
	if r < 0x80 {
		return append(out, byte(r))
	}
	if r < 0x800 {
		return append(out, byte(0xC0|r>>6), byte(0x80|r&0x3F))
	}
	if r < 0x10000 {
		return append(out, byte(0xE0|r>>12), byte(0x80|(r>>6)&0x3F), byte(0x80|r&0x3F))
	}
	return append(out,
		byte(0xF0|r>>18), byte(0x80|(r>>12)&0x3F),
		byte(0x80|(r>>6)&0x3F), byte(0x80|r&0x3F))
}
