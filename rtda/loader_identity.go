package rtda

import "sync/atomic"

// LoaderIdentity is an opaque, comparable identity for a defining class loader.
// Each class loader receives a unique *LoaderIdentity at creation time.
// The zero value (nil) is not a valid identity — callers must use VMIdentity
// or a loader-allocated identity.
//
// Identity comparison uses pointer equality: two LoaderIdentity values are
// equal iff they were allocated by the same NewLoaderIdentity call or are
// both VMIdentity.
type LoaderIdentity struct {
	id uint64
}

// VMIdentity is the canonical identity used by primitive types and void.
// These types are not defined by any class loader — they belong to the VM.
var VMIdentity = &LoaderIdentity{id: 0}

var loaderIDSeq uint64 = 1

// NewLoaderIdentity allocates a fresh, unique loader identity.
func NewLoaderIdentity() *LoaderIdentity {
	return &LoaderIdentity{id: atomic.AddUint64(&loaderIDSeq, 1)}
}
