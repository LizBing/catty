# ADR-0013: Direct Go runtime integration (no JVM abstraction layer)

- **Status:** Proposed
- **Date:** 2026-07-12

## Context

Traditional JVM I/O, networking, and threading sit above the JVM abstraction
layer (`java.io` → JNI → OS). catty's AOT path can directly map Java standard
library calls to Go runtime primitives — skipping the JVM abstraction layer
entirely. This isn't "bridging" — it's **compile-time replacement** of Java API
calls with equivalent Go calls.

## Decision

In the AOT path, Java standard library calls compile directly to Go runtime
calls.

- `Socket.connect()` → `net.Dial()`.
- `ServerSocket.accept()` → `net.Listener.Accept()`.
- `Thread.start()` → `go func()`.
- `Object.wait()` → `sync.Cond.Wait()`.
- `System.gc()` → `runtime.GC()`.
- `Thread.sleep()` → `time.Sleep()`.
- `FileInputStream.read()` → `os.File.Read()`.
- NIO `Selector` → Go netpoll (epoll/kqueue).

AOT code's I/O path is identical to a hand-written Go program. The interpreter
path still uses traditional JNI-like bridging.

## Consequences

**Positive**
- I/O performance equals native Go programs (no JVM abstraction overhead).
- Java NIO `Selector` uses Go's netpoll directly — no JNI, no epoll wrapper.
- `Thread.sleep(1000)` compiles to `time.Sleep(time.Second)` — one function call,
  not `Thread.sleep` → native → OS sleep.
- Deployable as a single static binary with no JRE dependency.

**Negative**
- AOT must know which Java APIs map to which Go functions (pattern matching at
  transpile time). The mapping table must be maintained as the JDK evolves.
- The interpreter path (for dynamically-loaded classes) still needs traditional
  bridging, creating a dual-path maintenance burden.
