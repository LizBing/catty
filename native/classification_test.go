package native

import (
	"testing"
)

// TestClassifiedRegistryInvariants verifies the R2-A native classification
// contract. Every registered native must be classified; category checks must
// be consistent.
func TestClassifiedRegistryInvariants(t *testing.T) {
	seen := make(map[string]bool)
	for _, e := range classifiedRegistry {
		key := e.ClassName + ":" + e.MethodName + ":" + e.Descriptor

		// Check 1: no duplicate classifications.
		if seen[key] {
			t.Errorf("duplicate classification: %s", key)
		}
		seen[key] = true

		// Check 2: classification is valid.
		if e.Classification > CategoryUnsupported {
			t.Errorf("%s: invalid classification %d", key, e.Classification)
		}

		// Check 3: Unsupported should not have a registered implementation.
		fn := GetNative(e.ClassName, e.MethodName, e.Descriptor)
		if e.Classification == CategoryUnsupported && fn != nil {
			t.Errorf("%s: classified Unsupported but has a registered implementation", key)
		}

		// Check 4: non-Unsupported must have a registered implementation.
		if e.Classification != CategoryUnsupported && fn == nil {
			t.Errorf("%s: classified %s but no registered implementation found", key, e.Classification)
		}
	}

	// Check 5: every registered native appears in the classified registry.
	nativeMu.RLock()
	defer nativeMu.RUnlock()
	for classAndMethod, fn := range nativeRegistry {
		// Parse key: className\x00methodName\x00descriptor
		parts := splitNulKey(classAndMethod)
		if len(parts) != 3 {
			t.Errorf("malformed registry key: %q", classAndMethod)
			continue
		}
		classKey := parts[0] + ":" + parts[1] + ":" + parts[2]
		if !seen[classKey] {
			_ = fn
			t.Errorf("registered native not classified: %s", classKey)
		}
	}
}

func splitNulKey(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
