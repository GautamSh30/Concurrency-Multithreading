# Go Concurrency Cheatsheet

## Goroutines
```go
go functionName(args)           // Launch goroutine
go func() { /* ... */ }()       // Anonymous goroutine
```

## Channels
```go
ch := make(chan T)              // Unbuffered (synchronous)
ch := make(chan T, n)           // Buffered (capacity n)
ch <- value                     // Send
v := <-ch                       // Receive
v, ok := <-ch                   // Receive with close check
close(ch)                       // Close channel
for v := range ch { }           // Iterate until closed
```

## Channel Direction
```go
func send(ch chan<- T) { }      // Send-only
func recv(ch <-chan T) { }      // Receive-only
```

## Select
```go
select {
case v := <-ch1:                // Receive from ch1
case ch2 <- val:                // Send to ch2
case <-time.After(d):           // Timeout
default:                        // Non-blocking
}
```

## WaitGroup
```go
var wg sync.WaitGroup
wg.Add(n)                       // Add n goroutines
go func() {
    defer wg.Done()             // Mark done
    // work...
}()
wg.Wait()                       // Block until all done
```

## Mutex
```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
// critical section

var rw sync.RWMutex
rw.RLock()                      // Multiple readers OK
rw.RUnlock()
```

## Atomic
```go
atomic.AddInt64(&counter, 1)
atomic.LoadInt64(&counter)
atomic.StoreInt64(&counter, val)
```

## Context
```go
ctx := context.Background()                    // Root context
ctx, cancel := context.WithCancel(parent)      // Cancellable
ctx, cancel := context.WithTimeout(parent, d)  // Timeout
ctx, cancel := context.WithDeadline(parent, t) // Deadline
ctx = context.WithValue(parent, key, val)      // Data bag
defer cancel()                                  // Always defer cancel!

<-ctx.Done()                                   // Wait for cancellation
ctx.Err()                                      // Why cancelled?
ctx.Value(key)                                 // Get value
```

## sync.Cond
```go
c := sync.NewCond(&sync.Mutex{})
c.L.Lock()
for !condition {
    c.Wait()          // Suspend goroutine
}
c.L.Unlock()

c.Signal()            // Wake one waiter
c.Broadcast()         // Wake all waiters
```

## sync.Once
```go
var once sync.Once
once.Do(func() {
    // runs only once, even from multiple goroutines
})
```

## sync.Pool
```go
pool := sync.Pool{
    New: func() interface{} { return new(bytes.Buffer) },
}
b := pool.Get().(*bytes.Buffer)
defer pool.Put(b)
```

## Race Detector
```bash
go run -race main.go
go test -race ./...
go build -race
```

## Pipeline Pattern
```go
func stage(in <-chan T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for v := range in {
            out <- transform(v)
        }
    }()
    return out
}
```

## Cancellation Pattern
```go
func worker(done <-chan struct{}, in <-chan T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for v := range in {
            select {
            case out <- process(v):
            case <-done:
                return
            }
        }
    }()
    return out
}
// To cancel: close(done)
```

## Common Pitfalls
| Pitfall | Fix |
|---------|-----|
| Goroutine leak | Use done channel or context cancellation |
| Loop variable capture | Pass as argument: `go func(v T) { }(val)` |
| Nil channel read/write | Always initialize with make() |
| Close channel twice | Only owner should close |
| Race condition | Use -race flag, mutex, or channels |
