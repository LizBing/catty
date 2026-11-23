package verify

import (
	"errors"
	"strings"
	"testing"

	"catty/internal/classfile"
	"catty/rt"
)

func TestCopyStateRefusesBeforeCopy(t *testing.T) {
	ctx := &Context{Work: &Work{Limit: 0}}
	s := &frameState{
		locals: []vtype{{kind: vtInt}, {kind: vtRef, class: "java/lang/Object"}},
		stack:  []vtype{{kind: vtInt}},
	}
	got, err := copyState(ctx, s)
	if err == nil || got != nil {
		t.Fatalf("refused copy returned state %#v, err %v", got, err)
	}
	var workErr *WorkError
	if !errors.As(err, &workErr) || workErr.Operation != "frame state copy" {
		t.Fatalf("expected a frame state copy refusal, got %T %v", err, err)
	}
	if ctx.Work.Used != 0 {
		t.Fatalf("refused copy consumed %d units", ctx.Work.Used)
	}
}

func TestStatesEqualRefusesBeforeComparison(t *testing.T) {
	ctx := &Context{Work: &Work{Limit: 0}}
	a := &frameState{locals: []vtype{{kind: vtInt}}}
	b := &frameState{locals: []vtype{{kind: vtInt}}}
	equal, err := statesEqual(ctx, a, b)
	if err == nil || equal {
		t.Fatalf("refused comparison returned %v, err %v", equal, err)
	}
	var workErr *WorkError
	if !errors.As(err, &workErr) || workErr.Operation != "frame state comparison" {
		t.Fatalf("expected a frame state comparison refusal, got %T %v", err, err)
	}
	if ctx.Work.Used != 0 {
		t.Fatalf("refused comparison consumed %d units", ctx.Work.Used)
	}
}

func TestHierarchyTraversalStopsAtRefusedVisit(t *testing.T) {
	ctx := &Context{
		Classes: map[string]*classfile.Class{
			"A":  {Name: "A", Interfaces: []string{"I1"}},
			"I1": {Name: "I1", Flags: rt.AccInterface | rt.AccAbstract, Interfaces: []string{"I2"}},
			"I2": {Name: "I2", Flags: rt.AccInterface | rt.AccAbstract, Interfaces: []string{"T"}},
			"T":  {Name: "T", Flags: rt.AccInterface | rt.AccAbstract},
		},
		Work: &Work{Limit: 2},
	}
	found, err := classImplements(ctx, "A", "T")
	if err == nil || found {
		t.Fatalf("refused traversal returned %v, err %v", found, err)
	}
	var workErr *WorkError
	if !errors.As(err, &workErr) || workErr.Operation != "interface hierarchy visit" {
		t.Fatalf("expected an interface hierarchy refusal, got %T %v", err, err)
	}
	if ctx.Work.Used != 2 {
		t.Fatalf("refused traversal consumed %d units, want exactly 2", ctx.Work.Used)
	}
	// The refused third visit must not have been performed: raising the budget
	// by one reaches the target through the same path.
	ctx.Work = &Work{Limit: 3}
	found, err = classImplements(ctx, "A", "T")
	if err != nil || !found {
		t.Fatalf("traversal with complete budget returned %v, err %v", found, err)
	}
}

func TestVerifierStopsAtEachRefusedStepWithoutPublishingState(t *testing.T) {
	c := handlerFixture([]classfile.VType{{Kind: classfile.VTInteger}})
	full := &Work{Limit: 1 << 20}
	successContext := &Context{Classes: map[string]*classfile.Class{}, Work: full}
	if err := VerifyClass(c, successContext); err != nil {
		t.Fatal(err)
	}
	if full.Used == 0 {
		t.Fatal("verifier charged no work")
	}
	seen := map[string]bool{}
	frameCopyLimit := int64(-1)
	for limit := int64(1); limit < full.Used; limit++ {
		work := &Work{Limit: limit}
		ctx := &Context{Classes: map[string]*classfile.Class{}, Work: work}
		err := VerifyClass(c, ctx)
		if err == nil {
			t.Fatalf("limit %d accepted but the full run used %d", limit, full.Used)
		}
		var workErr *WorkError
		if !errors.As(err, &workErr) {
			t.Fatalf("limit %d: expected WorkError, got %T %v", limit, err, err)
		}
		seen[workErr.Operation] = true
		if work.Used > limit {
			t.Fatalf("limit %d: refusal consumed %d units", limit, work.Used)
		}
		if len(ctx.Heights) != 0 || len(ctx.Shapes) != 0 {
			t.Fatalf("limit %d: refused verification published analysis state", limit)
		}
		if workErr.Operation == "frame state copy" && frameCopyLimit < 0 {
			frameCopyLimit = limit
		}
	}
	for _, want := range []string{
		"frame state copy",
		"frame state comparison",
		"handler range inspection",
		"exception edge",
		"instruction",
		"control-flow edge",
	} {
		if !seen[want] {
			t.Fatalf("no low-budget run failed at %q; saw %v", want, seen)
		}
	}
	if frameCopyLimit < 0 {
		t.Fatal("no successful charge boundary refused the frame copy step")
	}
	// At the exact copy boundary the refusal happens before the copy: the
	// verifier reports a LIMIT diagnostic and never records heights.
	work := &Work{Limit: frameCopyLimit}
	ctx := &Context{Classes: map[string]*classfile.Class{}, Work: work}
	err := VerifyClass(c, ctx)
	var workErr *WorkError
	if !errors.As(err, &workErr) || workErr.Operation != "frame state copy" {
		t.Fatalf("copy boundary: %T %v", err, err)
	}
	if work.Used != frameCopyLimit {
		t.Fatalf("copy boundary used %d of %d units", work.Used, frameCopyLimit)
	}
	if strings.Contains(err.Error(), "frame state copied") {
		t.Fatalf("copy refusal reported completed work: %v", err)
	}
}
