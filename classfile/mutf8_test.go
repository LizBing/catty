package classfile

import (
	"reflect"
	"testing"
)

func TestDecodeMUTF8ToUTF16ASCII(t *testing.T) {
	b := []byte("Hello")
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{'H', 'e', 'l', 'l', 'o'}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ASCII: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16NUL(t *testing.T) {
	// MUTF-8 encodes NUL (U+0000) as 0xC0 0x80 (overlong 2-byte form).
	b := []byte{0xC0, 0x80}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NUL: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16NULBetweenASCII(t *testing.T) {
	// Test NUL encoded between ASCII characters.
	b := []byte{'A', 0xC0, 0x80, 'B'}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{'A', 0, 'B'}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NUL between ASCII: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16Supplementary(t *testing.T) {
	// MUTF-8 encodes supplementary characters as two 3-byte surrogate halves.
	// U+10400 (𐐀) → high surrogate U+D801, low surrogate U+DC00
	// High: 0xED 0xA0 0x81 (0xD801 encoded as 3-byte)
	// Low:  0xED 0xB0 0x80 (0xDC00 encoded as 3-byte)
	b := []byte{0xED, 0xA0, 0x81, 0xED, 0xB0, 0x80}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0xD801, 0xDC00}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Supplementary: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16LoneHighSurrogate(t *testing.T) {
	// A lone high surrogate (0xD801) encoded as 3-byte MUTF-8, without a
	// following low surrogate. Must be preserved as-is.
	b := []byte{0xED, 0xA0, 0x81}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0xD801}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Lone high surrogate: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16LoneLowSurrogate(t *testing.T) {
	b := []byte{0xED, 0xB0, 0x80}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0xDC00}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Lone low surrogate: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16TwoByte(t *testing.T) {
	// 2-byte MUTF-8: U+00E9 (é) → 0xC3 0xA9
	b := []byte{0xC3, 0xA9}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0x00E9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("2-byte: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16ThreeByte(t *testing.T) {
	// 3-byte MUTF-8: U+4E2D (中) → 0xE4 0xB8 0xAD
	b := []byte{0xE4, 0xB8, 0xAD}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{0x4E2D}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("3-byte: got %v, want %v", got, want)
	}
}

func TestDecodeMUTF8ToUTF16Empty(t *testing.T) {
	got := decodeMUTF8ToUTF16([]byte{})
	if len(got) != 0 {
		t.Fatalf("Empty: expected empty, got %v", got)
	}
}

func TestDecodeMUTF8ToUTF16Mixed(t *testing.T) {
	// Mix of ASCII, NUL, supplementary, and lone surrogate.
	// A, NUL (0xC0 0x80), B, 中 (0xE4 0xB8 0xAD), lone high (0xED 0xA0 0x81)
	b := []byte{'A', 0xC0, 0x80, 'B', 0xE4, 0xB8, 0xAD, 0xED, 0xA0, 0x81}
	got := decodeMUTF8ToUTF16(b)
	want := []uint16{'A', 0, 'B', 0x4E2D, 0xD801}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Mixed: got %v, want %v", got, want)
	}
}

func TestMUTF8RoundTrip(t *testing.T) {
	// decodeMUTF8 (to Go string) and decodeMUTF8ToUTF16 (to []uint16) should
	// agree on ASCII.
	b := []byte("HelloWorld")
	s := decodeMUTF8(b)
	u := decodeMUTF8ToUTF16(b)
	for i, c := range s {
		if uint16(c) != u[i] {
			t.Fatalf("Round-trip mismatch at %d: string=%d utf16=%d", i, c, u[i])
		}
	}
}

func TestMUTF8RoundTripSupplementary(t *testing.T) {
	// Supplementary characters are recombined in decodeMUTF8 but kept as pairs
	// in decodeMUTF8ToUTF16.
	b := []byte{0xED, 0xA0, 0x81, 0xED, 0xB0, 0x80}
	s := decodeMUTF8(b) // should be a single rune U+10400
	u := decodeMUTF8ToUTF16(b) // should be two units 0xD801 0xDC00
	if len([]rune(s)) != 1 {
		t.Fatalf("decodeMUTF8 supplementary: expected 1 rune, got %d", len([]rune(s)))
	}
	if len(u) != 2 {
		t.Fatalf("decodeMUTF8ToUTF16 supplementary: expected 2 units, got %d", len(u))
	}
	if u[0] != 0xD801 || u[1] != 0xDC00 {
		t.Fatalf("decodeMUTF8ToUTF16 supplementary: got [%x, %x]", u[0], u[1])
	}
}
