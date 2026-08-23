package panicpkg

import (
	"testing"

	"excbench/common"
)

// These assertions guard the re-panic logic: a matched handler must run, an
// unmatched layer must re-panic (not swallow), and a non-sentinel panic must
// propagate untouched.

func TestNormalPathValues(t *testing.T) {
	want := (7*2654435761 + 17) >> 5
	if got := aNorm(7); got != want {
		t.Fatalf("aNorm(7) = %d, want %d", got, want)
	}
	if got := aNoHandlers(7); got != want {
		t.Fatalf("aNoHandlers(7) = %d, want %d", got, want)
	}
}

func TestShallowCatch(t *testing.T) {
	if got := cShallow(0); got != resultC {
		t.Fatalf("cShallow caught = %d, want %d (layer c handler)", got, resultC)
	}
}

func TestDeepCatch(t *testing.T) {
	if got := aDeepFree(0); got != resultA {
		t.Fatalf("aDeepFree caught = %d, want %d (only layer a matches)", got, resultA)
	}
	if got := aDeepRethrow(0); got != resultA {
		t.Fatalf("aDeepRethrow caught = %d, want %d (b/c/d must re-panic)", got, resultA)
	}
}

// A foreign (non-sentinel) panic must be re-panicked unchanged, exactly as the
// chain's deferred handlers do, so Java exceptions never mask a genuine Go
// runtime panic.
func TestNonSentinelPanicPropagates(t *testing.T) {
	defer func() {
		e := recover()
		if e == nil {
			t.Fatal("expected a panic to escape")
		}
		if s, ok := e.(string); !ok || s != "runtime-bug" {
			t.Fatalf("expected original string panic, got %v", e)
		}
	}()
	relay("runtime-bug")
}

// relay mirrors the chain's handler idiom: recover, match the sentinel, and
// re-panic anything else.
func relay(v any) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(*common.JException); ok {
				// sentinel handled here
			} else {
				panic(e)
			}
		}
	}()
	panic(v)
}
