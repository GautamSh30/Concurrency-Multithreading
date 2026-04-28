# Section 11: Pipelines, Fan-out & Fan-in, Cancellation

## What are Pipelines?
- Series of stages connected by channels
- Each stage: takes data in → performs operation → sends data out
- Process streams or batches of data
- Efficient use of I/O and multiple CPU cores
- Each stage is a goroutine

```
G1 ──ch1──> G2 ──ch2──> G3 ──ch3──> G4
(generator)  (stage1)    (stage2)   (consumer)
```

## Pipeline Stage Properties
- Takes a common `done` channel + input channel
- Returns an output channel  
- A stage can consume and return the same type (composability)

## Fan-out
- Multiple goroutines read from single channel
- Distribute work to parallelize CPU/I/O usage
- Useful for computationally intensive stages

## Fan-in (Merge)
- Combine results from multiple channels into one
- Create merge goroutines that read from multiple inputs

```
        ┌──> G2a ──ch2a──┐
G1 ─ch1─┤──> G2b ──ch2b──┼─Merge─ch3─> G3
        └──> G2c ──ch2c──┘
        (fan-out)         (fan-in)
```

## Goroutine Leak Problem
- If receiver only needs subset of values and abandons channel
- Upstream goroutines block forever on send → LEAK

## Cancellation Pattern
- Pass read-only `done` channel to all goroutines
- Close `done` to broadcast signal to all goroutines
- Use `select` to make send/receive preemptible

## Pipeline Construction Guidelines
1. Stages close outbound channels when all sends done
2. Stages receive from inbound until closed OR senders unblocked
3. Pipelines unblock senders by explicitly signaling when receiver abandons

## Image Processing Pipeline
- G1: Generate list of images to process
- G2: Generate thumbnail images
- G3: Store thumbnail images to disk
