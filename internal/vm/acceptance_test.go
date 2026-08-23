package vm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"catty/internal/kernel"
)

// runFixture loads testdata/<name> and runs its static main with an empty
// String[] argument, returning captured stdout.
func runFixture(t *testing.T, name string) (string, error) {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	cls, err := k.LoadClassBytes(data)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	th := New(k)
	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, err := k.NewArray("Ljava/lang/String;", 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = th.Call(main, nil, []kernel.Value{argsArr})
	return out.String(), err
}

func TestAcceptanceHelloWorld(t *testing.T) {
	out, err := runFixture(t, "HelloWorld.class")
	if err != nil {
		var th *kernel.Thrown
		if errors.As(err, &th) {
			t.Fatalf("uncaught %v", err)
		}
		t.Fatal(err)
	}
	if out != "Hello, Catty!\n" {
		t.Errorf("stdout = %q, want %q", out, "Hello, Catty!\n")
	}
}

func TestAcceptanceCollectionsDemo(t *testing.T) {
	out, err := runFixture(t, "CollectionsDemo.class")
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	want := strings.Join([]string{
		"150",
		"caught: Index: 99, Size: 5",
		"div-by-zero caught",
		"sum=150",
		"has30",
		"",
	}, "\n")
	if out != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out, want)
	}
}

func TestAcceptanceClassPathMulti(t *testing.T) {
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)

	cls, err := loader.Load("app/Main")
	if err != nil {
		t.Fatalf("load app/Main: %v", err)
	}
	th := New(k)
	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
	if _, err := th.Call(main, nil, []kernel.Value{argsArr}); err != nil {
		t.Fatalf("run: %v", err)
	}

	want := strings.Join([]string{
		"hi:m1:1",
		"hi:m2:2",
		"true",
		"0", // second Helper instance: fresh field default
		"",
	}, "\n")
	if out.String() != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out.String(), want)
	}
}

// TestAcceptanceThreads runs the multithreaded fixture through start/join,
// monitor contention and interrupt-of-sleep, comparing to the reference
// JVM output captured at fixture creation time.
func TestAcceptanceThreads(t *testing.T) {
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)

	cls, err := loader.Load("ThreadsDemo")
	if err != nil {
		t.Fatalf("load ThreadsDemo: %v", err)
	}
	th := New(k)
	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
	if _, err := th.Call(main, nil, []kernel.Value{argsArr}); err != nil {
		t.Fatalf("run: %v", err)
	}

	want := strings.Join([]string{
		"500",
		"interrupted ok",
		"main done main",
		"w1 alive=false",
		"",
	}, "\n")
	if out.String() != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out.String(), want)
	}
}

// syncBuf is a mutex-guarded buffer: Java threads println from their own
// goroutines while the test goroutine polls stdout.
type syncBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (w *syncBuf) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *syncBuf) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.String()
}

// TestAcceptanceHttpEcho boots the pure-Java HTTP echo server on an
// ephemeral port (args wiring), then talks real sockets to it.
func TestAcceptanceHttpEcho(t *testing.T) {
	out := &syncBuf{}
	k := kernel.New(kernel.Options{Stdout: out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)

	cls, err := loader.Load("HttpEcho")
	if err != nil {
		t.Fatalf("load HttpEcho: %v", err)
	}
	th := New(k)
	mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatal(err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 1)
	argsArr.Elems[0] = k.InternGo("0") // port 0 → ephemeral
	go th.Call(mainM, nil, []kernel.Value{argsArr})

	// wait for "listening <port>"
	dl := time.Now().Add(5 * time.Second)
	var port int
	for {
		line := strings.TrimSpace(out.String())
		if strings.HasPrefix(line, "listening ") {
			port, err = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "listening ")))
			if err != nil {
				t.Fatalf("bad listening line %q", line)
			}
			break
		}
		if time.Now().After(dl) {
			t.Fatalf("server never listened; got %q", out.String())
		}
		time.Sleep(5 * time.Millisecond)
	}

	req := "GET /abc HTTP/1.1\r\nHost: test\r\nX-A: b\r\n\r\n"
	conn, derr := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if derr != nil {
		t.Fatalf("dial: %v", derr)
	}
	conn.Write([]byte(req))
	resp, _ := io.ReadAll(conn)
	conn.Close()

	respStr := string(resp)
	if !strings.HasPrefix(respStr, "HTTP/1.0 200 OK\r\n") {
		t.Fatalf("status line: %q", respStr)
	}
	wantBody := strconv.Itoa(strings.Index(req, "\r\n\r\n"))
	if !strings.HasSuffix(respStr, "\r\n\r\n"+wantBody) {
		t.Fatalf("body mismatch: %q (want suffix %q)", respStr, wantBody)
	}
}
