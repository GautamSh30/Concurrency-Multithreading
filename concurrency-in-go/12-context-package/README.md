# Section 12: Context Package

## Two Primary Purposes
1. **Data bag**: Transport request-scoped data through call-graph
2. **Cancellation**: Cancel branches of call-graph

## Context is Safe
- Safe for simultaneous use by multiple goroutines
- Single context can be passed to any number of goroutines
- Cancelling context signals ALL goroutines to abandon work

## Creating Context
- `context.Background()`: Empty root context, never cancelled, used by main
- `context.TODO()`: Placeholder when unsure which context to use

## Cancellation Functions
- `context.WithCancel(parent)` → ctx, cancel: cancel() closes done channel
- `context.WithDeadline(parent, time)` → ctx, cancel: closes when clock passes deadline
- `context.WithTimeout(parent, duration)` → ctx, cancel: closes after duration (wrapper over WithDeadline)
- cancel() does not wait for work to stop
- cancel() is idempotent (safe to call multiple times)

## WithTimeout vs WithDeadline
- WithTimeout: countdown starts from context creation
- WithDeadline: set explicit clock time for expiry

## Data Bag
- `context.WithValue(parent, key, val)`: Associate request-scoped data
- `ctx.Value(key)`: Extract value by key
- Use only for request-scoped data, NOT essential function parameters

## Go Idioms for Context
1. Incoming requests should create a Context
2. Outgoing calls should accept a Context
3. Pass Context to any function performing I/O (as first parameter)
4. Any change creates new Context value propagated forward
5. When parent cancelled → all children cancelled
6. Use TODO if unsure which Context to use
7. Use context values only for request-scoped data

## HTTP Server Timeouts
- Important to conserve resources and protect from DDOS
- Four timeouts: ReadTimeout, WriteTimeout, ReadHeaderTimeout, IdleTimeout
- `http.TimeoutHandler()`: Returns 503 if handler exceeds time limit
- Request.Context() propagates cancellation down call graph
