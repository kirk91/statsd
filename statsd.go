package statsd

import (
	"os"
	"strings"
	"time"
)

type options struct {
	timeout       time.Duration
	flushPeriod   time.Duration
	maxPacketSize int
	errHandler    func(error)

	prefix   string
	hostname string
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

type Client struct {
	opts options

	cc *clientConn
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

	cc, err := newClientConn(network, addr, c)
	if err != nil {
		return nil, err
	}

	c.cc = cc

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

func (c *Client) Timingf(duration time.Duration, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), template, args))
}

func (c *Client) TimingSincef(start time.Time, template string, args ...interface{}) {
	c.send(c.encodeTpl(MetricTypeTiming, Float64(float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond)), template, args))
}

func (c *Client) IncrementWithHostname(bucket ...Field) {
	c.CountInt32WithHostname(1, bucket...)
}

func (c *Client) CountInt32WithHostname(n int32, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeCount, Int32(n), bucket))
}

func (c *Client) CountUint32WithHostname(n uint32, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeCount, Uint32(n), bucket))
}

func (c *Client) CountInt64WithHostname(n int64, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeCount, Int64(n), bucket))
}

func (c *Client) CountUint64WithHostname(n uint64, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeCount, Uint64(n), bucket))
}

func (c *Client) GaugeInt32WithHostname(n int32, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeGauge, Int32(n), bucket))
}

func (c *Client) GaugeUint32WithHostname(n uint32, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeGauge, Uint32(n), bucket))
}

func (c *Client) GaugeInt64WithHostname(n int64, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeGauge, Int64(n), bucket))
}

func (c *Client) GaugeUint64WithHostname(n uint64, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeGauge, Uint64(n), bucket))
}

func (c *Client) TimingSinceWithHostname(start time.Time, bucket ...Field) {
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / float64(time.Millisecond)
	c.send(c.encodeWithHostname(MetricTypeTiming, Float64(elapsed), bucket))
}

func (c *Client) TimingWithHostname(duration time.Duration, bucket ...Field) {
	c.send(c.encodeWithHostname(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), bucket))
}

func (c *Client) IncrementfWithHostname(template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeCount, Int32(1), template, args))
}

func (c *Client) CountInt32fWithHostname(n int32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeCount, Int32(n), template, args))
}

func (c *Client) CountUint32fWithHostname(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeCount, Uint32(n), template, args))
}

func (c *Client) CountInt64fWithHostname(n int64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeCount, Int64(n), template, args))
}

func (c *Client) CountUint64fWithHostname(n uint64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeCount, Uint64(n), template, args))
}

func (c *Client) GaugeInt32fWithHostname(n int32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeGauge, Int32(n), template, args))
}

func (c *Client) GaugeUint32fWithHostname(n uint32, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeGauge, Uint32(n), template, args))
}

func (c *Client) GaugeInt64fWithHostname(n int64, template string, args ...interface{}) {
	c.send(c.encodeTplWithHostname(MetricTypeGauge, Int64(n), template, args))
}

func (c *Client) GaugeUint64fWithHostname(n uint64, template string, args ...interface{}) {
	b := c.encodeTplWithHostname(MetricTypeGauge, Uint64(n), template, args)
	c.send(b)
}

func (c *Client) TimingfWithHostname(duration time.Duration, template string, args ...interface{}) {
	b := c.encodeTplWithHostname(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), template, args)
	c.send(b)
}

func (c *Client) TimingSincefWithHostname(start time.Time, template string, args ...interface{}) {
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / float64(time.Millisecond)
	b := c.encodeTplWithHostname(MetricTypeTiming, Float64(elapsed), template, args)
	c.send(b)
}

func (c *Client) encode(typ MetricType, val Field, bucket []Field) *buf {
	return encode(typ, val, c.opts.prefix, "", bucket)
}

func (c *Client) encodeWithHostname(typ MetricType, val Field, bucket []Field) *buf {
	return encode(typ, val, c.opts.prefix, c.opts.hostname, bucket)
}

func (c *Client) encodeTpl(typ MetricType, val Field, template string, fmtArgs []interface{}) *buf {
	return encodeTpl(typ, val, c.opts.prefix, "", template, fmtArgs)
}

func (c *Client) encodeTplWithHostname(typ MetricType, val Field, template string, fmtArgs []interface{}) *buf {
	return encodeTpl(typ, val, c.opts.prefix, c.opts.hostname, template, fmtArgs)
}

func (c *Client) send(b *buf) {
	if b == nil {
		return
	}
	c.cc.write(b.Bytes())
	freeBuf(b)
}
