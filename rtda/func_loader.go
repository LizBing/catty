package rtda

// FuncLoader is a minimal Loader implementation backed by a function.
// It is intended for tests and simple scenarios. Each call to LoadClass
// delegates to the provided function.
type FuncLoader struct {
	ID       *LoaderIdentity
	LoadFn   func(name string) *Class
	LoadResFn func(name string) ClassLoadResult
}

// NewFuncLoader creates a FuncLoader. If id is nil, a fresh identity is allocated.
func NewFuncLoader(loadFn func(name string) *Class) *FuncLoader {
	return &FuncLoader{
		ID:     NewLoaderIdentity(),
		LoadFn: loadFn,
	}
}

func (l *FuncLoader) LoadClass(name string) *Class {
	if l.LoadFn == nil {
		panic("catty: FuncLoader.LoadClass: no LoadFn configured")
	}
	return l.LoadFn(name)
}

func (l *FuncLoader) LoadClassResult(name string) ClassLoadResult {
	if l.LoadResFn != nil {
		return l.LoadResFn(name)
	}
	// Fallback: delegate to LoadClass and catch panics.
	defer func() {
		recover()
	}()
	c := l.LoadClass(name)
	if c != nil {
		return NewClassResult(c)
	}
	return NewFailureResult(&ClassLoadFailure{
		Kind: FailureNotFound,
		Name: name,
	})
}

func (l *FuncLoader) LoaderIdentity() *LoaderIdentity {
	if l.ID == nil {
		l.ID = NewLoaderIdentity()
	}
	return l.ID
}
