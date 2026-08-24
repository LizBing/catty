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
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"catty/internal/gen"
	"catty/internal/interop"
	"catty/internal/kernel"
	"catty/internal/vm"
)

func main() {
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
	k := kernel.New(kernel.Options{Stdout: os.Stdout})
	th := vm.New(k)

	err = interop.Bind(k, interop.Spec{
		Class: "go/Demo",
		Funcs: map[string]any{
			"md5Hex": func(s string) string { sum := md5.Sum([]byte(s)); return hex.EncodeToString(sum[:]) },
			"fetchLen": func(url string) (int64, error) {
				resp, err := http.Get(url)
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

	main, err := k.ResolveMethod(cls, "run", "()V")
	if err != nil {
		fmt.Fprintln(os.Stderr, "resolve run:", err)
		os.Exit(1)
	}
	if _, err := th.Call(main, nil, nil); err != nil {
		if tt, ok := err.(*kernel.Thrown); ok {
			fmt.Fprint(os.Stderr, kernel.FormatUncaught("main", tt))
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "engine error:", err)
		os.Exit(70)
	}
	fmt.Println("embeddemo: done")
}
