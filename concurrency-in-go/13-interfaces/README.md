# Section 13: Interfaces

## What are Interfaces?
- Abstract types that define behaviours (set of method signatures)
- "If something can do this, then it can be used here"
- Satisfied **implicitly** - no `implements` keyword

## Interface Values
- Contain **dynamic type** and **dynamic value**
- Method call through interface uses dynamic dispatch
- Compiler generates code to get method from type descriptor

## Purpose of Interfaces
- **Encapsulation**: Logic within user-defined types
- **Abstraction**: Higher level functions with behavior guarantees
- **Decoupling**: Definition decoupled from implementation

## Implicit Satisfaction
- Types just need to possess required methods
- Definition and implementation can be in different packages
- Can define new interfaces satisfied by existing types without changing them
- "Definition of interface is decoupled from implementation"

## Conventions
- Define interface when 2+ concrete types need uniform treatment
- Create smaller interfaces with fewer methods
- "Ask only for what is needed"

## Interfaces from Standard Library
### io.Writer
- One of the most widely used interfaces
- Abstraction for types to which bytes can be written: files, buffers, network, HTTP
```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

### fmt.Stringer
- Controls how values are printed
```go
type Stringer interface {
    String() string
}
```

## Interface Satisfaction
- Type satisfies interface if it implements ALL required methods
- `*os.File` satisfies Reader, Writer, Closer, ReadWriter
- `*bytes.Buffer` satisfies Reader, Writer - NOT Closer (no Close method)
- **Assignability**: Expression assignable to interface only if type satisfies it
- Interface wraps concrete type - only interface methods revealed

## Type Assertion
- Extract dynamic value from interface
- `v := x.(T)` - panics if wrong type
- `v, ok := x.(T)` - safe version (ok=false, v=zero value on failure)

## Type Switch
- Discover dynamic type of interface variable
```go
switch v := i.(type) {
case int: ...
case string: ...
default: ...
}
```

## Empty Interface
- `interface{}` (or `any` in Go 1.18+) specifies no methods
- Can hold any value
- CAUTION: Loses static type safety, use sparingly
- Prefer specific types or interfaces with required methods
