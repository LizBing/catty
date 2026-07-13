package rtda

import (
	"testing"
)

func TestStringValueNewAndLen(t *testing.T) {
	sv := NewStringValue([]uint16{'H', 'e', 'l', 'l', 'o'})
	if sv.Len() != 5 {
		t.Fatalf("Len: expected 5, got %d", sv.Len())
	}
	if sv.IsEmpty() {
		t.Fatal("IsEmpty: expected false for non-empty string")
	}
}

func TestStringValueEmpty(t *testing.T) {
	sv := NewStringValue([]uint16{})
	if sv.Len() != 0 {
		t.Fatalf("Len: expected 0, got %d", sv.Len())
	}
	if !sv.IsEmpty() {
		t.Fatal("IsEmpty: expected true for empty string")
	}
}

func TestStringValueCharAt(t *testing.T) {
	sv := NewStringValue([]uint16{0x41, 0x42, 0xD800, 0xDC00})
	if sv.CharAt(0) != 0x41 {
		t.Fatalf("CharAt(0): expected 0x41, got 0x%x", sv.CharAt(0))
	}
	if sv.CharAt(2) != 0xD800 {
		t.Fatalf("CharAt(2): expected 0xD800, got 0x%x", sv.CharAt(2))
	}
	if sv.CharAt(3) != 0xDC00 {
		t.Fatalf("CharAt(3): expected 0xDC00, got 0x%x", sv.CharAt(3))
	}
}

func TestStringValueDefensiveCopy(t *testing.T) {
	orig := []uint16{1, 2, 3}
	sv := NewStringValue(orig)
	orig[0] = 99 // mutate original slice
	if sv.CharAt(0) != 1 {
		t.Fatal("NewStringValue did not defensively copy the input slice")
	}
}

func TestStringValueUnitsDefensive(t *testing.T) {
	sv := NewStringValue([]uint16{1, 2, 3})
	u := sv.Units()
	u[0] = 99
	if sv.CharAt(0) != 1 {
		t.Fatal("Units() did not defensively copy")
	}
}

func TestStringValueHashCode(t *testing.T) {
	// Java String.hashCode for "Hello" (UTF-16 code units: H=72, e=101, l=108, l=108, o=111)
	// h = ((72*31+101)*31+108)*31+108)*31+111 = 69609650
	sv := NewStringValue([]uint16{'H', 'e', 'l', 'l', 'o'})
	expected := int32(69609650)
	if sv.HashCode() != expected {
		t.Fatalf("HashCode: expected %d, got %d", expected, sv.HashCode())
	}
}

func TestStringValueHashCodeEmpty(t *testing.T) {
	sv := NewStringValue([]uint16{})
	if sv.HashCode() != 0 {
		t.Fatalf("HashCode empty: expected 0, got %d", sv.HashCode())
	}
}

func TestStringValueHashCodeSupplementary(t *testing.T) {
	// UTF-16 for U+10400 (𐐀): high surrogate 0xD801, low surrogate 0xDC00
	// hashCode = 0xD801 * 31 + 0xDC00 = 55297*31 + 56320 = 1714207 + 56320 = 1770527
	sv := NewStringValue([]uint16{0xD801, 0xDC00})
	expected := int32(0xD801)*31 + int32(0xDC00)
	if sv.HashCode() != expected {
		t.Fatalf("HashCode supplementary: expected %d, got %d", expected, sv.HashCode())
	}
}

func TestStringValueEquals(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b', 'c'})
	b := NewStringValue([]uint16{'a', 'b', 'c'})
	if !a.Equals(b) {
		t.Fatal("Equals: identical content should be equal")
	}
}

func TestStringValueEqualsDifferent(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b', 'c'})
	b := NewStringValue([]uint16{'a', 'b', 'd'})
	if a.Equals(b) {
		t.Fatal("Equals: different content should not be equal")
	}
}

func TestStringValueEqualsDifferentLength(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b'})
	b := NewStringValue([]uint16{'a', 'b', 'c'})
	if a.Equals(b) {
		t.Fatal("Equals: different lengths should not be equal")
	}
}

func TestStringValueCompareTo(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b', 'c'})
	b := NewStringValue([]uint16{'a', 'b', 'd'})
	if a.CompareTo(b) >= 0 {
		t.Fatal("CompareTo: 'abc' < 'abd'")
	}
	if b.CompareTo(a) <= 0 {
		t.Fatal("CompareTo: 'abd' > 'abc'")
	}
}

func TestStringValueCompareToPrefix(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b'})
	b := NewStringValue([]uint16{'a', 'b', 'c'})
	if a.CompareTo(b) >= 0 {
		t.Fatal("CompareTo: 'ab' < 'abc'")
	}
}

func TestStringValueCompareToEqual(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b'})
	b := NewStringValue([]uint16{'a', 'b'})
	if a.CompareTo(b) != 0 {
		t.Fatal("CompareTo: equal strings should return 0")
	}
}

func TestStringValueSubstring(t *testing.T) {
	sv := NewStringValue([]uint16{'H', 'e', 'l', 'l', 'o'})
	sub := sv.Substring(1, 4) // "ell"
	if sub.Len() != 3 {
		t.Fatalf("Substring Len: expected 3, got %d", sub.Len())
	}
	if sub.CharAt(0) != 'e' || sub.CharAt(1) != 'l' || sub.CharAt(2) != 'l' {
		t.Fatal("Substring: unexpected content")
	}
}

func TestStringValueSubstringDefensive(t *testing.T) {
	sv := NewStringValue([]uint16{'H', 'e', 'l', 'l', 'o'})
	sub := sv.Substring(1, 4) // "ell"
	// Mutating the parent or child should not affect the other.
	if sv.CharAt(0) != 'H' {
		t.Fatal("Substring: parent was mutated")
	}
	if sub.CharAt(1) != 'l' {
		t.Fatal("Substring: child was mutated")
	}
}

func TestStringValueConcat(t *testing.T) {
	a := NewStringValue([]uint16{'a', 'b'})
	b := NewStringValue([]uint16{'c', 'd'})
	c := a.Concat(b)
	if c.Len() != 4 {
		t.Fatalf("Concat Len: expected 4, got %d", c.Len())
	}
	if c.CharAt(0) != 'a' || c.CharAt(1) != 'b' || c.CharAt(2) != 'c' || c.CharAt(3) != 'd' {
		t.Fatal("Concat: unexpected content")
	}
}

func TestStringValueIndexOf(t *testing.T) {
	sv := NewStringValue([]uint16{'a', 'b', 'c', 'a', 'b', 'c'})
	if sv.IndexOf('b') != 1 {
		t.Fatalf("IndexOf 'b': expected 1, got %d", sv.IndexOf('b'))
	}
	if sv.IndexOf('z') != -1 {
		t.Fatalf("IndexOf 'z': expected -1, got %d", sv.IndexOf('z'))
	}
}

func TestStringValueIndexOfSupplementary(t *testing.T) {
	// Surrogate pair for U+10400 (𐐀) = 0xD801 0xDC00
	sv := NewStringValue([]uint16{'a', 0xD801, 0xDC00, 'b'})
	if sv.IndexOf(0x10400) != 1 {
		t.Fatalf("IndexOf supplementary: expected 1, got %d", sv.IndexOf(0x10400))
	}
}

func TestStringValueIndexOfOutOfRange(t *testing.T) {
	sv := NewStringValue([]uint16{'a'})
	if sv.IndexOf(-1) != -1 {
		t.Fatal("IndexOf -1: expected -1")
	}
	if sv.IndexOf(0x110000) != -1 {
		t.Fatal("IndexOf 0x110000: expected -1")
	}
}

func TestStringValueStartsWith(t *testing.T) {
	sv := NewStringValue([]uint16{'a', 'b', 'c', 'd'})
	prefix := NewStringValue([]uint16{'a', 'b'})
	if !sv.StartsWith(prefix) {
		t.Fatal("StartsWith: 'abcd' should start with 'ab'")
	}
	notPrefix := NewStringValue([]uint16{'b', 'c'})
	if sv.StartsWith(notPrefix) {
		t.Fatal("StartsWith: 'abcd' should not start with 'bc'")
	}
}

func TestStringValueEndsWith(t *testing.T) {
	sv := NewStringValue([]uint16{'a', 'b', 'c', 'd'})
	suffix := NewStringValue([]uint16{'c', 'd'})
	if !sv.EndsWith(suffix) {
		t.Fatal("EndsWith: 'abcd' should end with 'cd'")
	}
	notSuffix := NewStringValue([]uint16{'b', 'c'})
	if sv.EndsWith(notSuffix) {
		t.Fatal("EndsWith: 'abcd' should not end with 'bc'")
	}
}

func TestStringValueToCharArray(t *testing.T) {
	sv := NewStringValue([]uint16{'a', 'b', 'c'})
	arr := sv.ToCharArray()
	if len(arr) != 3 || arr[0] != 'a' {
		t.Fatal("ToCharArray: unexpected content")
	}
	arr[0] = 'x'
	if sv.CharAt(0) != 'a' {
		t.Fatal("ToCharArray: defensive copy failed")
	}
}

func TestStringValueGoStringASCII(t *testing.T) {
	sv := NewStringValue([]uint16{'H', 'e', 'l', 'l', 'o'})
	if sv.GoString() != "Hello" {
		t.Fatalf("GoString: expected 'Hello', got %q", sv.GoString())
	}
}

func TestStringValueGoStringSupplementary(t *testing.T) {
	// U+10400 (𐐀) = surrogate pair 0xD801 0xDC00
	sv := NewStringValue([]uint16{0xD801, 0xDC00})
	s := sv.GoString()
	if len([]rune(s)) != 1 || []rune(s)[0] != 0x10400 {
		t.Fatalf("GoString supplementary: expected single rune U+10400, got %q", s)
	}
}

func TestStringValueGoStringLoneSurrogate(t *testing.T) {
	// Lone high surrogate → '?'
	sv := NewStringValue([]uint16{'A', 0xD800, 'B'})
	s := sv.GoString()
	if s != "A?B" {
		t.Fatalf("GoString lone high surrogate: expected 'A?B', got %q", s)
	}
}

func TestStringValueGoStringLoneLowSurrogate(t *testing.T) {
	// Lone low surrogate → '?'
	sv := NewStringValue([]uint16{'A', 0xDC00, 'B'})
	s := sv.GoString()
	if s != "A?B" {
		t.Fatalf("GoString lone low surrogate: expected 'A?B', got %q", s)
	}
}

func TestStringValueGoStringEmpty(t *testing.T) {
	sv := NewStringValue([]uint16{})
	if sv.GoString() != "" {
		t.Fatalf("GoString empty: expected '', got %q", sv.GoString())
	}
}

func TestNewStringValueNoAlias(t *testing.T) {
	// NewStringValue must always defensively copy — no alias to the input.
	units := []uint16{1, 2, 3}
	sv := NewStringValue(units)
	if sv.Len() != 3 {
		t.Fatal("NewStringValue: unexpected length")
	}
	// Mutation through the original slice must not affect the StringValue.
	units[0] = 99
	if sv.CharAt(0) != 1 {
		t.Fatal("NewStringValue must not alias its input")
	}
}
