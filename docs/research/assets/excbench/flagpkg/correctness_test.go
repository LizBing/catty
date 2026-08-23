package flagpkg

import "testing"

func TestNormalPathValues(t *testing.T) {
	want := (7*2654435761 + 17) >> 5
	v, err := aNorm(7)
	if err != nil {
		t.Fatalf("aNorm(7) err = %v, want nil", err)
	}
	if v != want {
		t.Fatalf("aNorm(7) = %d, want %d", v, want)
	}
}

func TestShallowCatch(t *testing.T) {
	v, err := cShallow(0)
	if err != nil {
		t.Fatalf("cShallow err = %v, want nil (handled)", err)
	}
	if v != resultC {
		t.Fatalf("cShallow caught = %d, want %d", v, resultC)
	}
}

func TestDeepCatch(t *testing.T) {
	v, err := aDeep(0)
	if err != nil {
		t.Fatalf("aDeep err = %v, want nil (handled)", err)
	}
	if v != resultA {
		t.Fatalf("aDeep caught = %d, want %d", v, resultA)
	}
}
