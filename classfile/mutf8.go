package classfile

// decodeMUTF8 decodes the "modified UTF-8" used by class files (JVMS §4.4.7) into
// a standard Go string. Modified UTF-8 differs from standard UTF-8 in two ways:
//   - the NUL code point (U+0000) is encoded as 0xC0 0x80 (an overlong 2-byte
//     form), never as a single 0x00 byte;
//   - supplementary characters (above U+FFFF) are encoded as a surrogate pair,
//     each half as a 3-byte sequence, rather than a single 4-byte sequence.
//
// For ASCII (the common case — class names and most descriptors) this is a
// zero-conversion pass. Malformed input falls back to the raw bytes.
func decodeMUTF8(b []byte) string {
	// Fast path: pure ASCII needs no conversion. Anything >= 0x80 forces the
	// slow path. This keeps the dominant case allocation-free beyond the string.
	ascii := true
	for _, c := range b {
		if c >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		return string(b)
	}

	var out []byte
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c < 0x80:
			out = append(out, c)
			i++
		case c&0xE0 == 0xC0: // 2-byte: 110xxxxx 10xxxxxx
			if i+1 >= len(b) {
				return string(b)
			}
			out = append(out, c, b[i+1])
			i += 2
		case c&0xF0 == 0xE0: // 3-byte: 1110xxxx 10xxxxxx 10xxxxxx
			if i+2 >= len(b) {
				return string(b)
			}
			// Detect a surrogate-pair half (0xED 0xA0..0xBF / 0xB0..0xBF) and
			// recombine two halves into one standard 4-byte UTF-8 sequence.
			if c == 0xED && i+5 < len(b) && b[i+1]&0xF0 == 0xA0 && b[i+3] == 0xED && b[i+4]&0xF0 == 0xB0 {
				hi := (int(c&0x0F) << 12) | (int(b[i+1]&0x3F) << 6) | int(b[i+2]&0x3F)
				lo := (int(b[i+3]&0x0F) << 12) | (int(b[i+4]&0x3F) << 6) | int(b[i+5]&0x3F)
				cp := 0x10000 + (((hi - 0xD800) << 10) | (lo - 0xDC00))
				out = append(out, byte(0xF0|(cp>>18)), byte(0x80|((cp>>12)&0x3F)),
					byte(0x80|((cp>>6)&0x3F)), byte(0x80|(cp&0x3F)))
				i += 6
				continue
			}
			out = append(out, c, b[i+1], b[i+2])
			i += 3
		default:
			out = append(out, c)
			i++
		}
	}
	return string(out)
}

// decodeMUTF8ToUTF16 decodes modified UTF-8 into UTF-16 code units losslessly.
// Unlike decodeMUTF8, which recombines surrogate pairs into 4-byte UTF-8 (and
// loses isolated surrogates), this function preserves every code unit exactly:
//   - NUL (U+0000) → single 0x0000 unit
//   - Valid surrogate pair → two UTF-16 code units (high then low)
//   - Lone high or low surrogate → the unit itself, preserved
//
// This is the accessor for CONSTANT_String entries in the constant pool.
// Names, descriptors, and other classfile text should continue to use
// decodeMUTF8, whose text semantics are sufficient for those purposes.
func decodeMUTF8ToUTF16(b []byte) []uint16 {
	if len(b) == 0 {
		return []uint16{}
	}

	// Fast path: pure ASCII.
	ascii := true
	for _, c := range b {
		if c >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		out := make([]uint16, len(b))
		for i, c := range b {
			out[i] = uint16(c)
		}
		return out
	}

	var out []uint16
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c < 0x80:
			out = append(out, uint16(c))
			i++
		case c >= 0xC0 && c < 0xE0: // 2-byte: 110xxxxx 10xxxxxx
			// Overlong NUL (0xC0 0x80) → U+0000.
			if c == 0xC0 && i+1 < len(b) && b[i+1] == 0x80 {
				out = append(out, 0)
				i += 2
				continue
			}
			if i+1 >= len(b) {
				for _, rb := range b {
					out = append(out, uint16(rb))
				}
				return out
			}
			cp := (int(c&0x1F) << 6) | int(b[i+1]&0x3F)
			out = append(out, uint16(cp))
			i += 2
		case c >= 0xE0 && c < 0xF0: // 3-byte: 1110xxxx 10xxxxxx 10xxxxxx
			if i+2 >= len(b) {
				for _, rb := range b {
					out = append(out, uint16(rb))
				}
				return out
			}
			cp := (int(c&0x0F) << 12) | (int(b[i+1]&0x3F) << 6) | int(b[i+2]&0x3F)
			// Surrogate-pair halves in MUTF-8: 0xED 0xA0..0xBF / 0xB0..0xBF
			// represent the high and low halves of a supplementary character.
			// Preserve them as-is (two UTF-16 code units).
			out = append(out, uint16(cp))
			i += 3
		default:
			// Malformed byte → preserve as code unit.
			out = append(out, uint16(c))
			i++
		}
	}
	return out
}
