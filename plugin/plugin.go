// Package plugin is the narrow launcher adapter that connects the Gate 3
// admission service to the runtime. Generated Gate 3 applications register
// Admit before installing their image. The runtime itself never imports
// internal/admission, which keeps the shared verifier free of a dependency
// cycle.
package plugin

import (
	"catty/internal/admission"
	"catty/rt"
)

// Admit validates and stages one plugin JAR candidate. The runtime supplies a
// detached parent-class snapshot so admission runs outside shared runtime
// locks.
func Admit(req *rt.AdmissionRequest) (*rt.AdmissionResult, error) {
	return admission.AdmitSnapshot(req.Bytes, req.Parent, req.ParentClasses, req.Profile)
}
