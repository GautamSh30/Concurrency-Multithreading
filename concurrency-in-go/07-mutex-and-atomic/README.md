# Section 7: Mutex & Atomic Operations

## When to Use Channels vs Mutex vs WaitGroup
| Tool | Use When |
|------|----------|
| **Channels** | Passing copy of data, distributing units of work, communicating async results |
| **Mutex** | Protecting caches, managing shared state |
| **WaitGroup** | Waiting for goroutines to finish, barrier/collection point |

## sync.Mutex
- Provides exclusive access to shared resource
- `mu.Lock()` / `mu.Unlock()` (use `defer mu.Unlock()`)
- Critical section = bottleneck between goroutines
- Developer convention - compiler doesn't enforce

## sync.RWMutex
- Multiple readers allowed simultaneously
- Writers get exclusive lock
- `mu.RLock()` / `mu.RUnlock()` for reads
- `mu.Lock()` / `mu.Unlock()` for writes

## sync/atomic
- Low-level atomic operations on memory
- Lockless - very fast for counters
- `atomic.AddUint64(&ops, 1)` / `atomic.LoadUint64(&ops)`
