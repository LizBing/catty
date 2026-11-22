package rt

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// trace is the opt-in JSON-lines evidence channel enabled by CATTY_TRACE.
// Java stdout remains untouched.
type trace struct {
	mu            sync.Mutex
	file          *os.File
	pid           int
	nextThrowable int64
}

func openTrace(path string) (*trace, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("catty: cannot open trace %s: %w", path, err)
	}
	return &trace{file: f, pid: os.Getpid()}, nil
}

func (tr *trace) close() {
	if tr == nil || tr.file == nil {
		return
	}
	_ = tr.file.Close()
	tr.file = nil
}

func (tr *trace) enabled() bool { return tr != nil && tr.file != nil }

func (t *Thread) trace() *trace {
	if t == nil || t.RT == nil {
		return nil
	}
	return t.RT.trace
}

func (tr *trace) emit(fields map[string]any) {
	if !tr.enabled() {
		return
	}
	fields["pid"] = tr.pid
	tr.mu.Lock()
	defer tr.mu.Unlock()
	data, err := json.Marshal(fields)
	if err != nil {
		return
	}
	data = append(data, '\n')
	_, _ = tr.file.Write(data)
}

func (tr *trace) throwableID(o *Object) int64 {
	if o == nil {
		return 0
	}
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if o.TraceID == 0 {
		tr.nextThrowable++
		o.TraceID = tr.nextThrowable
	}
	return o.TraceID
}

func methodFields(m *Method) map[string]any {
	return map[string]any{
		"owner":      m.Owner.Name,
		"name":       m.Name,
		"descriptor": m.Descriptor,
		"loaderId":   loaderID(m.Owner.Loader),
	}
}

func methodIDMap(m *Method) map[string]any { return methodFields(m) }

func (tr *trace) enter(thread int64, m *Method, depth int) {
	if !tr.enabled() {
		return
	}
	fields := methodFields(m)
	fields["event"] = "enter"
	fields["thread"] = thread
	fields["engine"] = m.Engine
	fields["depth"] = depth
	tr.emit(fields)
}

func (tr *trace) exit(thread int64, m *Method, depth int) {
	if !tr.enabled() {
		return
	}
	fields := methodFields(m)
	fields["event"] = "return"
	fields["thread"] = thread
	fields["depth"] = depth
	tr.emit(fields)
}

func (tr *trace) unwind(thread int64, m *Method, depth int, thrown *Object) {
	if !tr.enabled() {
		return
	}
	fields := methodFields(m)
	fields["event"] = "unwind"
	fields["thread"] = thread
	fields["depth"] = depth
	fields["throwable"] = tr.throwableID(thrown)
	fields["throwableLoaderId"] = loaderID(thrown.Class.Loader)
	tr.emit(fields)
}

func (tr *trace) dispatch(thread int64, opcode string, caller *Method, pc int32, declaration, target *Method, receiver *Class, hit bool) {
	if !tr.enabled() {
		return
	}
	callerFields := map[string]any{}
	if caller != nil {
		callerFields = methodIDMap(caller)
	}
	tr.emit(map[string]any{
		"event":            "dispatch",
		"thread":           thread,
		"opcode":           opcode,
		"caller":           callerFields,
		"pc":               pc,
		"declaration":      methodIDMap(declaration),
		"target":           methodIDMap(target),
		"receiverLoaderId": loaderID(receiver.Loader),
		"cache":            map[string]any{"hit": hit, "miss": !hit},
	})
}

func (tr *trace) reflect(thread int64, declaration, target *Method) {
	if !tr.enabled() {
		return
	}
	tr.emit(map[string]any{
		"event":       "reflect",
		"thread":      thread,
		"declaration": methodIDMap(declaration),
		"target":      methodIDMap(target),
	})
}

func (tr *trace) catchEvent(thread int64, m *Method, throwPC, handlerPC int32, thrown *Object) {
	if !tr.enabled() {
		return
	}
	fields := methodFields(m)
	fields["event"] = "catch"
	fields["thread"] = thread
	fields["throwPC"] = throwPC
	fields["handlerPC"] = handlerPC
	fields["throwable"] = tr.throwableID(thrown)
	fields["throwableLoaderId"] = loaderID(thrown.Class.Loader)
	tr.emit(fields)
}

func (tr *trace) wrap(thread int64, wrapper, cause *Object) {
	if !tr.enabled() {
		return
	}
	tr.emit(map[string]any{
		"event":           "wrap",
		"thread":          thread,
		"wrapper":         tr.throwableID(wrapper),
		"wrapperLoaderId": loaderID(wrapper.Class.Loader),
		"cause":           tr.throwableID(cause),
		"causeLoaderId":   loaderID(cause.Class.Loader),
	})
}

// admission emits one plugin_admission decision.
func (tr *trace) admission(thread int64, res *AdmissionResult, digest []byte, parent, loader *Loader) {
	if !tr.enabled() {
		return
	}
	fields := map[string]any{
		"event":              "plugin_admission",
		"thread":             thread,
		"result":             "accepted",
		"jarSHA256":          fmt.Sprintf("%x", digest),
		"parentLoaderId":     loaderID(parent),
		"loaderId":           loader.ID,
		"definedClasses":     len(res.Defined),
		"shadowedClasses":    len(res.Shadowed),
		"interpretedMethods": res.Interpreted,
		"deferredReferences": res.Deferred,
	}
	tr.emit(fields)
}

// admissionRejected records a refused registration. No loader ID is assigned.
func (tr *trace) admissionRejected(thread int64, pe *PluginError, digest []byte, parent *Loader) {
	if !tr.enabled() {
		return
	}
	code := "VERIFY"
	fields := map[string]any{
		"event":          "plugin_admission",
		"thread":         thread,
		"result":         "rejected",
		"code":           code,
		"parentLoaderId": loaderID(parent),
	}
	if pe != nil {
		fields["code"] = pe.Code
		if pe.Entry != "" {
			fields["entry"] = pe.Entry
		}
		if pe.Class != "" {
			fields["class"] = pe.Class
		}
		if pe.Method != "" {
			fields["method"] = pe.Method
		}
		if pe.Descriptor != "" {
			fields["descriptor"] = pe.Descriptor
		}
		if pe.PC != 0 {
			fields["pc"] = pe.PC
		}
		if pe.Feature != "" {
			fields["feature"] = pe.Feature
		}
		if pe.Limit != "" {
			fields["limit"] = pe.Limit
		}
		if pe.Message != "" {
			fields["message"] = pe.Message
		}
	}
	if len(digest) > 0 {
		fields["jarSHA256"] = fmt.Sprintf("%x", digest)
	}
	tr.emit(fields)
}

// classDefine records the first materialization of a plugin class.
func (tr *trace) classDefine(thread int64, c *Class, loader *Loader) {
	if !tr.enabled() {
		return
	}
	tr.emit(map[string]any{
		"event":     "class_define",
		"thread":    thread,
		"class":     c.Name,
		"loaderId":  loaderID(loader),
		"jarSHA256": loader.Digest,
	})
}

// classInit records initialization start and completion for one defining class.
func (tr *trace) classInit(thread int64, c *Class, start bool, thrown *Object) {
	if !tr.enabled() {
		return
	}
	event := "class_init"
	result := "complete"
	if start {
		result = "start"
	}
	fields := map[string]any{
		"event":    event,
		"thread":   thread,
		"class":    c.Name,
		"loaderId": loaderID(c.Loader),
		"phase":    result,
	}
	if thrown != nil {
		fields["failed"] = true
		fields["throwable"] = tr.throwableID(thrown)
	}
	tr.emit(fields)
}
