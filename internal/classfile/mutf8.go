package classfile

// decodeMUTF8 converts JVM modified UTF-8 (JVMS §4.4.7) to a Go string.
//
// Differences from standard UTF-8: NUL is encoded as 0xC0 0x80 and
// supplementary characters (U+10000..U+10FFFF) are encoded as CESU-8 style
// surrogate pairs of two 3-byte sequences.
func decodeMUTF8(b []byte) string {
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

	var runes []rune
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c&0x80 == 0: // 1-byte
			runes = append(runes, rune(c))
			i++
		case c == 0xC0: // 0xC0 0x80 → NUL
			if i+1 < len(b) && b[i+1] == 0x80 {
				runes = append(runes, 0)
				i += 2
			} else {
				runes = append(runes, rune(0xFFFD))
				i++
			}
		case c&0xE0 == 0xC0: // 2-byte
			if i+1 >= len(b) {
				runes = append(runes, rune(0xFFFD))
				i = len(b)
				break
			}
			runes = append(runes, rune(c&0x1F)<<6|rune(b[i+1]&0x3F))
			i += 2
		case c&0xF0 == 0xE0: // 3-byte, possibly high surrogate
			if i+2 >= len(b) {
				runes = append(runes, rune(0xFFFD))
				i = len(b)
				break
			}
			v := rune(c&0x0F)<<12 | rune(b[i+1]&0x3F)<<6 | rune(b[i+2]&0x3F)
			i += 3
			if v >= 0xD800 && v <= 0xDBFF && i+3 <= len(b) && b[i]&0xF0 == 0xE0 {
				v2 := rune(b[i]&0x0F)<<12 | rune(b[i+1]&0x3F)<<6 | rune(b[i+2]&0x3F)
				if v2 >= 0xDC00 && v2 <= 0xDFFF {
					runes = append(runes, ((v-0xD800)<<10|(v2-0xDC00))+0x10000)
					i += 3
					continue
				}
			}
			runes = append(runes, v)
		default: // invalid → replacement char, skip one byte
			runes = append(runes, rune(0xFFFD))
			i++
		}
	}
	return string(runes)
}

// EncodeMUTF8 converts a Go string back to JVM modified UTF-8.
// Used by tests and by runtime code that must synthesize constant pool data.
func EncodeMUTF8(s string) []byte {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		switch {
		case r == 0:
			out = append(out, 0xC0, 0x80)
		case r < 0x80:
			out = append(out, byte(r))
		case r < 0x800:
			out = append(out, byte(0xC0|r>>6), byte(0x80|r&0x3F))
		case r < 0x10000:
			out = append(out, byte(0xE0|r>>12), byte(0x80|(r>>6)&0x3F), byte(0x80|r&0x3F))
		default:
			u := r - 0x10000
			hi := rune(0xD800 + (u >> 10))
			lo := rune(0xDC00 + (u & 0x3FF))
			out = append(out,
				byte(0xE0|hi>>12), byte(0x80|(hi>>6)&0x3F), byte(0x80|hi&0x3F),
				byte(0xE0|lo>>12), byte(0x80|(lo>>6)&0x3F), byte(0x80|lo&0x3F))
		}
	}
	return out
}
