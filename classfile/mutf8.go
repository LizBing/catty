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
				hi := (int(b[i+1]&0x3F) << 6) | int(b[i+2]&0x3F)
				lo := (int(b[i+4]&0x3F) << 6) | int(b[i+5]&0x3F)
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
