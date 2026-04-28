# Section 8: Sync Primitives (Cond, Once, Pool)

## sync.Cond (Conditional Variable)
- Container of goroutines waiting for a certain condition
- Created with `sync.NewCond(&sync.Mutex{})` 
- Methods:
  - `c.Wait()`: Suspends goroutine, unlocks c.L, re-locks on wakeup. MUST be called in a loop!
  - `c.Signal()`: Wakes ONE goroutine (longest waiting)
  - `c.Broadcast()`: Wakes ALL waiting goroutines
- Use when multiple goroutines need to wait on multiple conditions

## sync.Once
- Runs a function exactly once, even across goroutines
- `once.Do(funcValue)` - only the first call executes the function
- Useful for one-time initialization (DB connections, config loading)

## sync.Pool
- Creates and makes available a pool of reusable objects
- `pool.Get()` retrieves an object, `pool.Put(obj)` returns it
- Reduces allocation overhead
- Objects may be garbage collected at any time
