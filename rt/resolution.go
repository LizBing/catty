package rt

// resolutionKey identifies one cached symbolic-resolution outcome. Java caches
// per referring runtime class and constant-pool entry; synthetic host
// operations use a nil Class and their member names instead.
type resolutionKey struct {
	Kind       byte
	Class      *Class
	CP         int
	Owner      string
	Name       string
	Descriptor string
}

// dispatchKey identifies one guarded compiled call site. The caller method and
// original bytecode offset are stable for the lifetime of the image.
type dispatchKey struct {
	Caller *Method
	PC     int32
	Op     uint8
}

// dispatchEntry is the single cached successful target for one call site.
type dispatchEntry struct {
	Declaration *Method
	Receiver    *Class
	Target      *Method
}

// AdmissionRequest is the read-only input the runtime hands to the registered
// admission callback. Bytes are a private copy owned by the callback.
// ParentClasses is the detached nearest-first snapshot of the parent chain's
// published definitions, so the adapter never reads live loader tables.
type AdmissionRequest struct {
	Bytes         []byte
	Parent        *Loader
	Profile       string
	ParentClasses map[string]*Class
}

// DeferredReference is one symbolic dependency retained for lazy runtime
// resolution because verification and structural checks succeeded without it.
type DeferredReference struct {
	Kind               string
	Owner              string
	Name               string
	Descriptor         string
	Referrer           string
	ReferrerMethod     string
	ReferrerDescriptor string
	PC                 int
}

// AdmissionResult carries staged immutable definitions back to the runtime.
type AdmissionResult struct {
	DeferredRefs []DeferredReference
	Classes      []*ClassDef
	Defined      []string
	Shadowed     []string
	Interpreted  int
	Deferred     int
	Digest       string
	ClassBytes   int64
	WorkUsed     int64
	Code         string
	Message      string
	Entry        string
	Class        string
	Method       string
	Descriptor   string
	PC           int
	Feature      string
	Limit        string
}
