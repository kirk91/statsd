# statsd
Fast, efficient statsd client

[![Build Status](https://travis-ci.org/kirk91/statsd.svg?branch=master)](https://travis-ci.org/kirk91/statsd)
[![GoDoc](https://godoc.org/github.com/kirk91/statsd?status.svg)](https://godoc.org/github.com/kirk91/statsd)


## Installation

```sh
go get -u github.com/kirk91/statsd
```
## Usage

```go
import "github.com/kirk91/statsd"

c, _ := statsd.New("udp", "127.0.0.1:8125")

// performance sensitive
c.Increment(statsd.String("foo"), statsd.String("bar"))
c.CountInt32(10, statsd.String("mong"), statsd.String("mew"))
c.GaugeInt32(1024, statsd.String("kong"), statsd.String("mew"))
c.Timing(time.Now(), statsd.String("kong"), statsd.Int(1))

// convenience sensitive
c.Incrementf("foo.bar")
c.CountInt32(10, "mong.new")
c.GaugeInt32(1024, "kong.new")
c.Timingf(time.Now(), "kong.1")
```

## Benchmark

- go1.10.3 darwin/amd64
- macos 10.13.6 (2.3 GHz Intel Core i5, 8 GB 1600 MHz DDR3)

```sh
❯ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/kirk91/statsd
BenchmarkIncrement-4           	10000000	       159 ns/op	       0 B/op	       0 allocs/op
BenchmarkIncrementParallel-4   	10000000	       135 ns/op	       0 B/op	       0 allocs/op
BenchmarkCount-4               	10000000	       159 ns/op	       0 B/op	       0 allocs/op
BenchmarkCountParallel-4       	10000000	       154 ns/op	       0 B/op	       0 allocs/op
BenchmarkGauge-4               	10000000	       157 ns/op	       0 B/op	       0 allocs/op
BenchmarkGaugeParallel-4       	10000000	       150 ns/op	       0 B/op	       0 allocs/op
BenchmarkTiming-4              	 5000000	       295 ns/op	       0 B/op	       0 allocs/op
BenchmarkTimingParallel-4      	 5000000	       268 ns/op	       0 B/op	       0 allocs/op
```

## Inspired by
- [statsd](https://github.com/alexcesaro/statsd)
- [zap](https://github.com/uber-go/zap)
