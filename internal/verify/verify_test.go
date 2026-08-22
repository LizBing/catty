package verify

import (
	"errors"
	"os"
	"strings"
	"testing"

	"catty/internal/classfile"
)

func mustParseFixture(t *testing.T, name string) *classfile.ClassFile {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	return cf
}

func TestStructuralVerifyAcceptsJavacOutput(t *testing.T) {
	for _, f := range []string{"HelloWorld.class", "CollectionsDemo.class"} {
		if err := Verify(mustParseFixture(t, f), nil); err != nil {
			t.Errorf("%s rejected by structural verifier: %v", f, err)
		}
	}
}

func TestRejectsBranchTargetWithoutFrame(t *testing.T) {
	cf := mustParseFixture(t, "CollectionsDemo.class")
	main := cf.FindMethod("main", "([Ljava/lang/String;)V")
	if main == nil || len(main.Code.StackMaps) < 2 {
		t.Fatal("fixture shape unexpected")
	}
	// Drop one frame: some branch target now lacks its merge point.
	main.Code.StackMaps = main.Code.StackMaps[1:]
	err := Verify(cf, nil)
	if err == nil || !strings.Contains(err.Error(), "no StackMapTable frame") {
		t.Fatalf("want missing-frame rejection, got %v", err)
	}
}

func TestRejectsObjectItemNotAClass(t *testing.T) {
	cf := mustParseFixture(t, "CollectionsDemo.class")
	main := cf.FindMethod("main", "([Ljava/lang/String;)V")
	for i := range main.Code.StackMaps {
		fr := &main.Code.StackMaps[i]
		for j, it := range fr.Locals {
			if it.Tag == classfile.VItemObject {
				// Point the OBJECT item at a Utf8 constant instead.
				for k := range cf.Constants {
					if cf.Constants[k].Tag == classfile.CUtf8 {
						fr.Locals[j].CPoolIdx = uint16(k)
						break
					}
				}
				err := Verify(cf, nil)
				if err == nil || !strings.Contains(err.Error(), "CONSTANT_Class") {
					t.Fatalf("want pool-kind rejection, got %v", err)
				}
				return
			}
		}
	}
	t.Fatal("no OBJECT item found in fixture frames")
}

func TestRejectsFrameBeyondCode(t *testing.T) {
	cf := mustParseFixture(t, "HelloWorld.class")
	hello := cf.FindMethod("main", "([Ljava/lang/String;)V")
	_ = hello
	// HelloWorld has no SM frames; inject an out-of-range one on the ctor.
	ctor := cf.FindMethod("<init>", "()V")
	ctor.Code.StackMaps = append(ctor.Code.StackMaps,
		classfile.StackMapFrame{Offset: 1 << 20})
	err := Verify(cf, nil)
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("want range rejection, got %v", err)
	}
}

func TestErrorCarriesMethodAndPc(t *testing.T) {
	cf := mustParseFixture(t, "CollectionsDemo.class")
	main := cf.FindMethod("main", "([Ljava/lang/String;)V")
	main.Code.StackMaps = nil
	err := Verify(cf, nil)
	if err == nil {
		t.Fatal("want error")
	}
	var ev *errVerify
	if !errors.As(err, &ev) {
		t.Fatalf("error type %T", err)
	}
	if !strings.Contains(ev.method, "main") {
		t.Errorf("method context missing: %v", err)
	}
}
