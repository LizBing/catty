package kernel

import "io"

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }

var _ io.Writer = discard{}
