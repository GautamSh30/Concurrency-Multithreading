# Section 10: Web Crawler - Sequential vs Concurrent

## Overview
- Build a web crawler using Go's concurrency features
- Start with sequential approach, then convert to concurrent
- Uses `golang.org/x/net/html` for HTML parsing
- Demonstrates real-world concurrency application

## Sequential Crawler
- Fetches one page at a time
- Simple but slow
- Easy to understand and debug

## Concurrent Crawler
- Fetches URLs in parallel using goroutines
- Uses channels to coordinate work
- Significantly faster
- Channel-based work distribution pattern
