// Command embeddemo demonstrates the interop productized surface
// (ADR-0011 / P-0010): a Go host embeds the Catty runtime, binds real Go
// capabilities (crypto md5, an HTTP fetch against a local server the host
// itself starts), and runs Java business rules that consume them.
//
// Run: go run ./cmd/embeddemo
package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"catty/internal/gen"
	"catty/internal/interop"
	"catty/internal/kernel"
	"catty/internal/vm"
)

func main() {
	requests := flag.Int("requests", 0, "sustained-load mode: total requests across workers (0 = single-shot demo)")
	conc := flag.Int("c", 4, "worker goroutines for load mode")
	flag.Parse()

	sharedClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 128,
			MaxConnsPerHost:     0,
		},
	}

	// 1. The host provides capabilities: a local HTTP service.
	mux := http.NewServeMux()
	mux.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"status":"ok","bytes":24,"source":"go-net-http"}`)
	})
	srv := &http.Server{Handler: mux}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	go srv.Serve(ln)
	defer srv.Close()
	base := fmt.Sprintf("http://%s", ln.Addr())

	// 2. Embed the runtime and bind Go functions into the Java namespace.
	// Load mode keeps rules silent; only the summary prints.
	var stdout io.Writer = os.Stdout
	if *requests > 0 {
		stdout = io.Discard
	}
	k := kernel.New(kernel.Options{Stdout: stdout})
	th := vm.New(k)

	err = interop.Bind(k, interop.Spec{
		Class: "go/Demo",
		Funcs: map[string]any{
			"md5Hex": func(s string) string { sum := md5.Sum([]byte(s)); return hex.EncodeToString(sum[:]) },
			"fetchLen": func(url string) (int64, error) {
				resp, err := sharedClient.Get(url)
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				n, err := io.Copy(io.Discard, resp.Body)
				return n, err
			},
			"baseUrl": func() string { return base },
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "bind:", err)
		os.Exit(1)
	}

	loader := kernel.NewClassPathLoader(k, []string{"testdata/cp"})
	k.SetResolver(loader.Load)
	cls, err := loader.Load("demo/Rules")
	if err != nil {
		fmt.Fprintln(os.Stderr, "load rules:", err)
		os.Exit(1)
	}
	gen.Install(k) // emitted rule bodies take precedence

	run, err := k.ResolveMethod(cls, "run", "()V")
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve run:", err)
		os.Exit(1)
	}

	if *requests > 0 {
		runLoad(k, run, *requests, *conc)
		return
	}

	if _, err := th.Call(run, nil, nil); err != nil {
		if tt, ok := err.(*kernel.Thrown); ok {
			fmt.Fprint(os.Stderr, k.FormatUncaught("main", tt))
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "engine error:", err)
		os.Exit(70)
	}
	fmt.Println("embeddemo: done")
}

// runLoad drives the Java rule as a sustained request stream from C worker
// goroutines, sampling host-side per-request latency, then prints
// percentiles — the embedding shape of selling points ②/③/④ combined.
func runLoad(k *kernel.Kernel, run *kernel.Method, total int, workers int) {
	th := make([]*vm.Thread, workers)
	for w := 0; w < workers; w++ {
		// One Java-thread context per worker: shared contexts would race
		// frame tracking and SOE metering.
		th[w] = vm.NewWithID(k, k.MintKey())
	}
	if workers < 1 {
		workers = 1
	}
	per := (total + workers - 1) / workers
	lat := make([][]time.Duration, workers)
	var done atomic.Int64
	var wg sync.WaitGroup
	start := time.Now()
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			samples := make([]time.Duration, 0, per)
			for i := 0; i < per; i++ {
				t0 := time.Now()
				if _, err := th[w].Call(run, nil, nil); err != nil {
					fmt.Fprintln(os.Stderr, "worker:", err)
					return
				}
				samples = append(samples, time.Since(t0))
			}
			lat[w] = samples
			done.Add(int64(len(samples)))
		}(w)
	}
	wg.Wait()
	wall := time.Since(start)

	all := make([]int64, 0, total)
	for _, s := range lat {
		for _, d := range s {
			all = append(all, int64(d))
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	n := len(all)
	pct := func(p float64) time.Duration {
		if n == 0 {
			return 0
		}
		idx := int(float64(n-1) * p)
		return time.Duration(all[idx])
	}
	fmt.Printf("requests=%d workers=%d wall=%s\n", n, workers, wall.Round(time.Millisecond))
	fmt.Printf("throughput=%.0f req/s\n", float64(n)/wall.Seconds())
	fmt.Printf("p50=%s p90=%s p99=%s max=%s\n",
		pct(0.50).Round(time.Microsecond), pct(0.90).Round(time.Microsecond),
		pct(0.99).Round(time.Microsecond),
		time.Duration(all[n-1]).Round(time.Microsecond))
}
