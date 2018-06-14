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
c.Increment(statsd.String("foo"), statsd.String("bar"))
c.Count(10, statsd.String("mong"), statsd.String("mew"))
c.GaugeInt(1024, statsd.String("kong"), statsd.String("mew"))
c.Timing(time.Now(), statsd.String("kong"), statsd.Int(1))
```

## Benchmark

- go1.8 darwin/amd64
- macos 10.12.3 (2.6 GHz Intel Core i5, 8 GB 1600 MHz DDR3)

```sh
BenchmarkIncrement-4    10000000               207 ns/op               0 B/op          0 allocs/op
BenchmarkCount-4        10000000               196 ns/op               0 B/op          0 allocs/op
BenchmarkGauge-4        10000000               210 ns/op               0 B/op          0 allocs/op
BenchmarkTiming-4       10000000               224 ns/op               0 B/op          0 allocs/op
```

## Inspired by
- [statsd](https://github.com/alexcesaro/statsd)
- [zap](https://github.com/uber-go/zap)

