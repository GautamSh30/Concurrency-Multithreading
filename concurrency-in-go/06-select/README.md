# Section 6: Select Statement

> Lectures 33–36 — Concurrency in Go (Deepak Kumar Gunjetti)

## What is Select?

- A `select` statement is like a `switch` where **each case specifies a channel operation** (send or receive).
- All channel operations are considered **simultaneously**.
- The statement **waits** until at least one case is ready to proceed, then executes that case.
- Enables **multiplexing** — reading from or writing to whichever channel is ready first.

```go
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case ch2 <- value:
    fmt.Println("sent to ch2")
}
```

---

## Key Properties

| Property | Behavior |
|---|---|
| Multiple cases ready | One is chosen **at random** (uniform pseudo-random selection) |
| `default` case present | Select becomes **non-blocking** — if no channel is ready, `default` runs immediately |
| Empty `select {}` | Blocks **forever** (useful to keep `main` alive) |
| All channels are `nil` | Blocks **forever** (nil channel operations always block) |

---

## Use Cases

### 1. Multiplexing

Read from whichever channel has data ready first:

```go
select {
case msg := <-ch1:
    fmt.Println("ch1:", msg)
case msg := <-ch2:
    fmt.Println("ch2:", msg)
}
```

### 2. Timeouts with `time.After()`

`time.After(d)` returns a `<-chan time.Time` that receives a value after duration `d`. Combine with `select` to implement timeouts:

```go
select {
case v := <-ch:
    fmt.Println("Received:", v)
case <-time.After(3 * time.Second):
    fmt.Println("Timeout!")
}
```

### 3. Non-blocking Communication with `default`

The `default` case makes `select` non-blocking — if no channel operation is ready, `default` executes immediately:

```go
select {
case msg := <-ch:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message available right now")
}
```

---

## Select in a Loop

A common pattern is to wrap `select` in a `for` loop to continuously multiplex:

```go
for {
    select {
    case msg := <-ch:
        fmt.Println(msg)
    case <-done:
        fmt.Println("shutting down")
        return
    }
}
```

---

## Nil Channel Trick

Setting a channel to `nil` inside a select effectively **disables** that case, since operations on nil channels block forever and `select` skips blocked cases:

```go
for ch1 != nil || ch2 != nil {
    select {
    case v, ok := <-ch1:
        if !ok {
            ch1 = nil
            continue
        }
        fmt.Println("ch1:", v)
    case v, ok := <-ch2:
        if !ok {
            ch2 = nil
            continue
        }
        fmt.Println("ch2:", v)
    }
}
```

---

## Key Takeaways

1. `select` is the control structure for working with **multiple channels** concurrently.
2. When multiple cases are ready, one is picked **at random** — no starvation, no priority.
3. Use `time.After()` with `select` for clean **timeout** handling.
4. Add a `default` case for **non-blocking** channel operations.
5. An empty `select {}` blocks forever — useful for keeping a program alive.
6. Setting a channel to `nil` disables its `select` case — a powerful pattern for merging multiple channels.
