package rt

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	goruntime "runtime"
)

// metricsRecorder is the opt-in Gate 5 process metrics sink. With metrics
// disabled the Runtime holds a nil recorder and every counting helper returns
// on the feature check. Snapshot-safe counter updates never depend on tracing
// and never take a runtime lock.
type metricsRecorder struct {
	file   *os.File
	mu     sync.Mutex
	closed bool
	pid    int

	classesLoaded    atomic.Int64
	methodsAOT       atomic.Int64
	methodsInterp    atomic.Int64
	methodsIntrinsic atomic.Int64
	methodsAbstract  atomic.Int64
	crossEngine      atomic.Int64
	reflective       atomic.Int64
	reflectionMisses atomic.Int64
	javaExceptions   atomic.Int64
	contentions      atomic.Int64

	liveJavaThreads atomic.Int64
	peakJavaThreads atomic.Int64

	postGCHeapBytes atomic.Uint64
	postGCHeapSeen  atomic.Bool
}

// SetMetricsPath opens the opt-in metrics output. The caller uses it before
// Install so class publication counters include bootstrap and application
// definitions; an open failure is a launcher setup error.
func (r *Runtime) SetMetricsPath(path string) error {
	if path == "" {
		return nil
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("catty: cannot open metrics %s: %w", path, err)
	}
	r.metrics = &metricsRecorder{file: f, pid: os.Getpid()}
	return nil
}

// CloseMetrics releases the metrics sink. It is safe after finish.
func (r *Runtime) CloseMetrics() {
	m := r.metrics
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.file != nil {
		_ = m.file.Close()
		m.file = nil
	}
	m.closed = true
}

func (m *metricsRecorder) classLoaded() {
	if m == nil {
		return
	}
	m.classesLoaded.Add(1)
}

func (m *metricsRecorder) methodPublished(engine string) {
	if m == nil {
		return
	}
	switch engine {
	case "aot":
		m.methodsAOT.Add(1)
	case "interp":
		m.methodsInterp.Add(1)
	case "intrinsic":
		m.methodsIntrinsic.Add(1)
	case "abstract":
		m.methodsAbstract.Add(1)
	}
}

func (m *metricsRecorder) crossEngineCall() {
	if m != nil {
		m.crossEngine.Add(1)
	}
}

func (m *metricsRecorder) reflectiveInvocation() {
	if m != nil {
		m.reflective.Add(1)
	}
}

func (m *metricsRecorder) reflectionMiss() {
	if m != nil {
		m.reflectionMisses.Add(1)
	}
}

func (m *metricsRecorder) javaException() {
	if m != nil {
		m.javaExceptions.Add(1)
	}
}

func (m *metricsRecorder) monitorContention() {
	if m != nil {
		m.contentions.Add(1)
	}
}

func (m *metricsRecorder) liveThreadStarted() {
	if m == nil {
		return
	}
	live := m.liveJavaThreads.Add(1)
	for {
		peak := m.peakJavaThreads.Load()
		if live <= peak || m.peakJavaThreads.CompareAndSwap(peak, live) {
			return
		}
	}
}

func (m *metricsRecorder) liveThreadEnded() {
	if m != nil {
		m.liveJavaThreads.Add(-1)
	}
}

// noteLiveThreadStartLocked publishes one started main, worker or hook peer to
// the live-thread high-water counter. The caller holds thrMu.
func (r *Runtime) noteLiveThreadStartLocked() {
	if r.metrics != nil {
		r.metrics.liveThreadStarted()
	}
}

// noteLiveThreadEndLocked publishes one terminating peer. The caller holds
// thrMu.
func (r *Runtime) noteLiveThreadEndLocked() {
	if r.metrics != nil {
		r.metrics.liveThreadEnded()
	}
}

// collectGarbage implements Runtime.gc: request a Go collection and, when
// metrics are enabled, record the process-wide Go heap observation immediately
// after the cycle. The observation is labeled as live Go-managed heap for the
// whole process, including runtime structures.
func (r *Runtime) collectGarbage() {
	goruntime.GC()
	if r.metrics == nil {
		return
	}
	var stats goruntime.MemStats
	goruntime.ReadMemStats(&stats)
	r.metrics.postGCHeapBytes.Store(stats.HeapAlloc)
	r.metrics.postGCHeapSeen.Store(true)
}

type metricsDocument struct {
	SchemaVersion  int             `json:"schemaVersion"`
	PID            int             `json:"pid"`
	Profile        string          `json:"profile"`
	TerminalReason string          `json:"terminalReason"`
	Status         int32           `json:"status"`
	Counters       metricsCounters `json:"counters"`
	GoRuntime      metricsGo       `json:"goRuntime"`
}

type metricsCounters struct {
	ClassesLoaded          int64            `json:"classesLoaded"`
	MethodsByEngine        map[string]int64 `json:"methodsByEngine"`
	CrossEngineCalls       int64            `json:"crossEngineCalls"`
	ReflectiveInvocations  int64            `json:"reflectiveInvocations"`
	ReflectionLookupMisses int64            `json:"reflectionLookupMisses"`
	JavaExceptions         int64            `json:"javaExceptions"`
	MonitorContentions     int64            `json:"monitorContentions"`
	LiveJavaThreads        int64            `json:"liveJavaThreads"`
	PeakJavaThreads        int64            `json:"peakJavaThreads"`
}

type metricsGo struct {
	HeapAllocBytes  uint64  `json:"heapAllocBytes"`
	TotalAllocBytes uint64  `json:"totalAllocBytes"`
	NumGC           uint32  `json:"numGC"`
	PostGCHeapBytes *uint64 `json:"postGCHeapBytes,omitempty"`
	PostGCHeapScope string  `json:"postGCHeapScope,omitempty"`
	HeapAllocScope  string  `json:"heapAllocScope"`
}

// finish writes the one terminal metrics document. A write failure is reported
// by the launcher on stderr and never changes the selected terminal status.
func (r *Runtime) finishMetrics(res terminalResult) error {
	m := r.metrics
	if m == nil {
		return nil
	}
	var stats goruntime.MemStats
	goruntime.ReadMemStats(&stats)
	counters := metricsCounters{
		ClassesLoaded:          m.classesLoaded.Load(),
		MethodsByEngine:        map[string]int64{},
		CrossEngineCalls:       m.crossEngine.Load(),
		ReflectiveInvocations:  m.reflective.Load(),
		ReflectionLookupMisses: m.reflectionMisses.Load(),
		JavaExceptions:         m.javaExceptions.Load(),
		MonitorContentions:     m.contentions.Load(),
		LiveJavaThreads:        m.liveJavaThreads.Load(),
		PeakJavaThreads:        m.peakJavaThreads.Load(),
	}
	for engine, count := range map[string]int64{
		"aot":       m.methodsAOT.Load(),
		"interp":    m.methodsInterp.Load(),
		"intrinsic": m.methodsIntrinsic.Load(),
		"abstract":  m.methodsAbstract.Load(),
	} {
		counters.MethodsByEngine[engine] = count
	}
	gostat := metricsGo{
		HeapAllocBytes:  stats.HeapAlloc,
		TotalAllocBytes: stats.TotalAlloc,
		NumGC:           stats.NumGC,
		HeapAllocScope:  "live Go-managed heap for the process, including runtime structures",
	}
	if m.postGCHeapSeen.Load() {
		value := m.postGCHeapBytes.Load()
		gostat.PostGCHeapBytes = &value
		gostat.PostGCHeapScope = "live Go-managed heap observed immediately after Runtime.gc()"
	}
	doc := metricsDocument{
		SchemaVersion:  1,
		PID:            m.pid,
		Profile:        r.Image.Profile,
		TerminalReason: res.reason,
		Status:         res.status,
		Counters:       counters,
		GoRuntime:      gostat,
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("catty: cannot encode metrics: %w", err)
	}
	data = append(data, '\n')
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || m.file == nil {
		return fmt.Errorf("catty: metrics output is closed")
	}
	if _, err := m.file.Write(data); err != nil {
		return fmt.Errorf("catty: cannot write metrics: %w", err)
	}
	return nil
}
