package runtime

import (
	"testing"

	"catty/rtda"
)

func TestInvokeTypedDynamicIsExplicitlyRejected(t *testing.T) {
	result := InvokeTypedDynamic(rtda.InvocationLookupRequest{})
	if !result.IsInternalFailure() || result.Failure() == nil {
		t.Fatal("typed dynamic AOT invocation must be explicitly rejected")
	}
}
