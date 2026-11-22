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
	tr.emit(fields)
}

func (tr *trace) dispatch(thread int64, opcode string, declaration, target *Method) {
	if !tr.enabled() {
		return
	}
	tr.emit(map[string]any{
		"event":       "dispatch",
		"thread":      thread,
		"opcode":      opcode,
		"declaration": methodIDMap(declaration),
		"target":      methodIDMap(target),
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
	tr.emit(fields)
}

func (tr *trace) wrap(thread int64, wrapper, cause *Object) {
	if !tr.enabled() {
		return
	}
	tr.emit(map[string]any{
		"event":   "wrap",
		"thread":  thread,
		"wrapper": tr.throwableID(wrapper),
		"cause":   tr.throwableID(cause),
	})
}
