package prototypes

import (
	"errors"
	"sync"
	"testing"
	"time"
)

var errAlreadyStarted = errors.New("already started")

type researchSupervisor struct {
	mu        sync.Mutex
	zero      *sync.Cond
	nonDaemon int
}

func newResearchSupervisor() *researchSupervisor {
	s := &researchSupervisor{}
	s.zero = sync.NewCond(&s.mu)
	return s
}

func (s *researchSupervisor) add(daemon bool) {
	if daemon {
		return
	}
	s.mu.Lock()
	s.nonDaemon++
	s.mu.Unlock()
}

func (s *researchSupervisor) done(daemon bool) {
	if daemon {
		return
	}
	s.mu.Lock()
	s.nonDaemon--
	if s.nonDaemon == 0 {
		s.zero.Broadcast()
	}
	s.mu.Unlock()
}

func (s *researchSupervisor) waitForNoNonDaemon() {
	s.mu.Lock()
	for s.nonDaemon != 0 {
		s.zero.Wait()
	}
	s.mu.Unlock()
}

type researchThread struct {
	mu         sync.Mutex
	started    bool
	terminated chan struct{}
	daemon     bool
	supervisor *researchSupervisor
}

func newResearchThread(supervisor *researchSupervisor, daemon bool) *researchThread {
	return &researchThread{
		terminated: make(chan struct{}),
		daemon:     daemon,
		supervisor: supervisor,
	}
}

func (t *researchThread) start(run func()) error {
	t.mu.Lock()
	if t.started {
		t.mu.Unlock()
		return errAlreadyStarted
	}
	t.started = true
	t.supervisor.add(t.daemon)
	go func() {
		defer func() {
			close(t.terminated)
			t.supervisor.done(t.daemon)
		}()
		run()
	}()
	t.mu.Unlock()
	return nil
}

func (t *researchThread) join() {
	t.mu.Lock()
	started := t.started
	t.mu.Unlock()
	if started {
		<-t.terminated
	}
}

func TestThreadStartOnceJoinAndNonDaemonLiveness(t *testing.T) {
	supervisor := newResearchSupervisor()
	thread := newResearchThread(supervisor, false)
	release := make(chan struct{})
	wrote := 0

	if err := thread.start(func() {
		<-release
		wrote = 42
	}); err != nil {
		t.Fatal(err)
	}
	if err := thread.start(func() {}); !errors.Is(err, errAlreadyStarted) {
		t.Fatalf("second start error=%v", err)
	}

	vmDone := make(chan struct{})
	go func() {
		supervisor.waitForNoNonDaemon()
		close(vmDone)
	}()
	select {
	case <-vmDone:
		t.Fatal("VM liveness ignored a running non-daemon thread")
	case <-time.After(time.Millisecond):
	}

	close(release)
	thread.join()
	if wrote != 42 {
		t.Fatalf("join did not publish worker write: %d", wrote)
	}
	select {
	case <-vmDone:
	case <-time.After(time.Second):
		t.Fatal("VM liveness did not observe non-daemon termination")
	}
}

func TestDaemonDoesNotKeepSupervisorAlive(t *testing.T) {
	supervisor := newResearchSupervisor()
	thread := newResearchThread(supervisor, true)
	release := make(chan struct{})
	if err := thread.start(func() { <-release }); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		supervisor.waitForNoNonDaemon()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("daemon thread incorrectly kept supervisor alive")
	}
	close(release)
	thread.join()
}
