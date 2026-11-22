package rt

import "math"

// Java integral conversion helpers. The JVM specifies saturating conversions
// from floating-point values to integral types: NaN becomes zero, values above
// the maximum saturate to the maximum and values below the minimum saturate to
// the minimum; in-range values round toward zero.

// F2I implements f2i.
func F2I(v float32) int32 {
	f := float64(v)
	switch {
	case math.IsNaN(f):
		return 0
	case f >= 2147483648.0:
		return math.MaxInt32
	case f <= -2147483648.0:
		return math.MinInt32
	default:
		return int32(f)
	}
}

// D2I implements d2i.
func D2I(v float64) int32 {
	switch {
	case math.IsNaN(v):
		return 0
	case v >= 2147483648.0:
		return math.MaxInt32
	case v <= -2147483648.0:
		return math.MinInt32
	default:
		return int32(v)
	}
}

// F2L implements f2l.
func F2L(v float32) int64 {
	f := float64(v)
	switch {
	case math.IsNaN(f):
		return 0
	case f >= 9223372036854775808.0:
		return math.MaxInt64
	case f <= -9223372036854775808.0:
		return math.MinInt64
	default:
		return int64(f)
	}
}

// D2L implements d2l.
func D2L(v float64) int64 {
	switch {
	case math.IsNaN(v):
		return 0
	case v >= 9223372036854775808.0:
		return math.MaxInt64
	case v <= -9223372036854775808.0:
		return math.MinInt64
	default:
		return int64(v)
	}
}

// I2B implements i2b.
func I2B(v int32) int32 { return int32(int8(v)) }

// I2C implements i2c.
func I2C(v int32) int32 { return int32(uint16(v)) }

// I2S implements i2s.
func I2S(v int32) int32 { return int32(int16(v)) }

// Lcmp implements lcmp.
func Lcmp(a, b int64) int32 {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Fcmp implements fcmpl and fcmpg. nanResult is -1 for fcmpl and 1 for fcmpg.
func Fcmp(a, b float32, nanResult int32) int32 {
	switch {
	case math.IsNaN(float64(a)) || math.IsNaN(float64(b)):
		return nanResult
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Dcmp implements dcmpl and dcmpg.
func Dcmp(a, b float64, nanResult int32) int32 {
	switch {
	case math.IsNaN(a) || math.IsNaN(b):
		return nanResult
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// ShiftDistance masks a Java shift distance for an int operand.
func IntShiftDistance(d int32) uint { return uint(uint32(d) & 31) }

// LongShiftDistance masks a Java shift distance for a long operand.
func LongShiftDistance(d int32) uint { return uint(uint32(d) & 63) }

// FloatToRawIntBits canonicalizes no NaN payload; Float.floatToIntBits uses
// this profile's raw-bit payload (the fixtures compare classification and
// defined bits only).
func FloatToIntBits(v float32) int32 {
	if math.IsNaN(float64(v)) {
		return int32(0x7fc00000)
	}
	return int32(math.Float32bits(v))
}

// DoubleToLongBits canonicalizes NaN like Double.doubleToLongBits.
func DoubleToLongBits(v float64) int64 {
	if math.IsNaN(v) {
		return int64(0x7ff8000000000000)
	}
	return int64(math.Float64bits(v))
}

// FRem implements frem.
func FRem(a, b float32) float32 { return float32(math.Mod(float64(a), float64(b))) }

// DRem implements drem.
func DRem(a, b float64) float64 { return math.Mod(a, b) }
