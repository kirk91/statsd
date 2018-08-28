package statsd

import (
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type options struct {
	prefix                       string
	hostname                     string
	timeout                      time.Duration
	flushPeriod                  time.Duration
	maxPacketSize                int
	disalbeMultiCoreOptimization bool

	errHandler func(error)
}

type Option func(*options)

func ErrorHandler(h func(error)) Option {
	return func(o *options) {
		o.errHandler = h
	}
}

func FlushPeriod(d time.Duration) Option {
	return func(o *options) {
		o.flushPeriod = d
	}
}

func Prefix(s string) Option {
	return func(o *options) {
		o.prefix = s
	}
}

func MaxPacketSize(n int) Option {
	return func(o *options) {
		o.maxPacketSize = n
	}
}

func Timeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

func Hostname(hostname string) Option {
	return func(o *options) {
		o.hostname = hostname
	}
}

func DisableMultiCoreOptimization() Option {
	return func(o *options) {
		o.disalbeMultiCoreOptimization = true
	}
}

type Client struct {
	opts options

	connCount int
	conns     []*clientConn

	metricCount uint64 // use to sharding metrics
}

func New(network, addr string, opt ...Option) (*Client, error) {
	c := &Client{}
	for _, o := range opt {
		o(&c.opts)
	}

	if c.opts.hostname == "" {
		hostname, _ := os.Hostname()
		c.opts.hostname = strings.Replace(hostname, ".", "_", -1)
	}
	if c.opts.timeout <= 0 {
		c.opts.timeout = time.Second * 5
	}
	if c.opts.flushPeriod <= 0 {
		c.opts.flushPeriod = time.Millisecond * 100
	}
	if c.opts.maxPacketSize <= 0 {
		c.opts.maxPacketSize = 1400
	}

	c.connCount = runtime.GOMAXPROCS(0)
	if c.opts.disalbeMultiCoreOptimization {
		c.connCount = 1
	}
	c.conns = make([]*clientConn, c.connCount)
	for i := 0; i < c.connCount; i++ {
		conn, err := newClientConn(network, addr, c)
		if err != nil {
			return nil, err
		}
		c.conns[i] = conn
	}

	return c, nil
}

func (c *Client) Increment(bucket ...Field) {
	c.CountInt32(1, bucket...)
}

func (c *Client) CountInt32(n int32, bucket ...Field) {
	c.send(c.encode(MetricTypeCount, Int32(n), bucket))
}

func (c *Client) CountUint32(n uint32, bucket ...Field) {
	c.send(c.encode(MetricTypeCount, Uint32(n), bucket))
}

func (c *Client) CountInt64(n int64, bucket ...Field) {
	c.send(c.encode(MetricTypeCount, Int64(n), bucket))
}

func (c *Client) CountUint64(n uint64, bucket ...Field) {
	c.send(c.encode(MetricTypeCount, Uint64(n), bucket))
}

func (c *Client) GaugeInt32(n int32, bucket ...Field) {
	c.send(c.encode(MetricTypeGauge, Int32(n), bucket))
}

func (c *Client) GaugeUint32(n uint32, bucket ...Field) {
	c.send(c.encode(MetricTypeGauge, Uint32(n), bucket))
}

func (c *Client) GaugeInt64(n int64, bucket ...Field) {
	c.send(c.encode(MetricTypeGauge, Int64(n), bucket))
}

func (c *Client) GaugeUint64(n uint64, bucket ...Field) {
	c.send(c.encode(MetricTypeGauge, Uint64(n), bucket))
}

func (c *Client) GaugeFloat64(n float64, bucket ...Field) {
	c.send(c.encode(MetricTypeGauge, Float64(n), bucket))
}

func (c *Client) TimingSince(start time.Time, bucket ...Field) {
	c.send(c.encode(MetricTypeTiming, Float64(float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond)), bucket))
}

func (c *Client) Timing(duration time.Duration, bucket ...Field) {
	c.send(c.encode(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), bucket))
}

func (c *Client) Incrementf(template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeCount, Int32(1), template, args))
}

func (c *Client) CountInt32f(n int32, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeCount, Int32(n), template, args))
}

func (c *Client) CountUint32f(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeCount, Uint32(n), template, args))
}

func (c *Client) CountInt64f(n int64, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeCount, Int64(n), template, args))
}

func (c *Client) CountUint64f(n uint64, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeCount, Uint64(n), template, args))
}

func (c *Client) GaugeInt32f(n int32, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeGauge, Int32(n), template, args))
}

func (c *Client) GaugeUint32f(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeGauge, Uint32(n), template, args))
}

func (c *Client) GaugeInt64f(n int64, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeGauge, Int64(n), template, args))
}

func (c *Client) GaugeUint64f(n uint64, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeGauge, Uint64(n), template, args))
}

func (c *Client) GaugeFloat64f(n float64, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeGauge, Float64(n), template, args))
}

func (c *Client) Timingf(duration time.Duration, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), template, args))
}

func (c *Client) TimingSincef(start time.Time, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeTiming, Float64(float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond)), template, args))
}

func (c *Client) IncrementWithHost(bucket ...Field) {
	c.CountInt32WithHost(1, bucket...)
}

func (c *Client) CountInt32WithHost(n int32, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeCount, Int32(n), bucket))
}

func (c *Client) CountUint32WithHost(n uint32, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeCount, Uint32(n), bucket))
}

func (c *Client) CountInt64WithHost(n int64, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeCount, Int64(n), bucket))
}

func (c *Client) CountUint64WithHost(n uint64, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeCount, Uint64(n), bucket))
}

func (c *Client) GaugeInt32WithHost(n int32, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeGauge, Int32(n), bucket))
}

func (c *Client) GaugeUint32WithHost(n uint32, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeGauge, Uint32(n), bucket))
}

func (c *Client) GaugeInt64WithHost(n int64, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeGauge, Int64(n), bucket))
}

func (c *Client) GaugeUint64WithHost(n uint64, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeGauge, Uint64(n), bucket))
}

func (c *Client) GaugeFloat64WithHost(n float64, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeGauge, Float64(n), bucket))
}

func (c *Client) TimingSinceWithHost(start time.Time, bucket ...Field) {
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / float64(time.Millisecond)
	c.send(c.encodeWithHost(MetricTypeTiming, Float64(elapsed), bucket))
}

func (c *Client) TimingWithHost(duration time.Duration, bucket ...Field) {
	c.send(c.encodeWithHost(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), bucket))
}

func (c *Client) IncrementfWithHost(template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeCount, Int32(1), template, args))
}

func (c *Client) CountInt32fWithHost(n int32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeCount, Int32(n), template, args))
}

func (c *Client) CountUint32fWithHost(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeCount, Uint32(n), template, args))
}

func (c *Client) CountInt64fWithHost(n int64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeCount, Int64(n), template, args))
}

func (c *Client) CountUint64fWithHost(n uint64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeCount, Uint64(n), template, args))
}

func (c *Client) GaugeInt32fWithHost(n int32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeGauge, Int32(n), template, args))
}

func (c *Client) GaugeUint32fWithHost(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeGauge, Uint32(n), template, args))
}

func (c *Client) GaugeInt64fWithHost(n int64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeGauge, Int64(n), template, args))
}

func (c *Client) GaugeUint64fWithHost(n uint64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeGauge, Uint64(n), template, args))
}

func (c *Client) GaugeFloat64fWithHost(n float64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHost(MetricTypeGauge, Float64(n), template, args))
}

func (c *Client) TimingfWithHost(duration time.Duration, template string, args ...interface{}) {
	b := c.encodeTplWithHost(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), template, args)
	c.send(b)
}

func (c *Client) TimingSincefWithHost(start time.Time, template string, args ...interface{}) {
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / float64(time.Millisecond)
	b := c.encodeTplWithHost(MetricTypeTiming, Float64(elapsed), template, args)
	c.send(b)
}

func (c *Client) encode(typ MetricType, val Field, bucket []Field) *buf {
	return encode(typ, val, c.opts.prefix, "", bucket)
}

func (c *Client) encodeWithHost(typ MetricType, val Field, bucket []Field) *buf {
	return encode(typ, val, c.opts.prefix, c.opts.hostname, bucket)
}

func (c *Client) encodeTpl(typ MetricType, val Field, template string, fmtArgs []interface{}) *buf {
	return encodeTpl(typ, val, c.opts.prefix, "", template, fmtArgs)
}

func (c *Client) encodeTplWithHost(typ MetricType, val Field, template string, fmtArgs []interface{}) *buf {
	return encodeTpl(typ, val, c.opts.prefix, c.opts.hostname, template, fmtArgs)
}

func (c *Client) send(b *buf) {
	if b == nil {
		return
	}

	idx := atomic.AddUint64(&c.metricCount, 1) % uint64(c.connCount)
	conn := c.conns[idx]
	conn.Write(b.Bytes())
	freeBuf(b)
}
