package statsd

import "time"

type options struct {
	timeout       time.Duration
	flushPeriod   time.Duration
	maxPacketSize int
	errHandler    func(error)

	prefix string
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

type Client struct {
	opts options

	cc *clientConn
}

func New(network, addr string, opt ...Option) (*Client, error) {
	c := &Client{}
	for _, o := range opt {
		o(&c.opts)
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
	c.send(encode(MetricTypeCount, Int32(n), c.opts.prefix, bucket))
}

func (c *Client) CountUint32(n uint32, bucket ...Field) {
	c.send(encode(MetricTypeCount, Uint32(n), c.opts.prefix, bucket))
}

func (c *Client) CountInt64(n int64, bucket ...Field) {
	c.send(encode(MetricTypeCount, Int64(n), c.opts.prefix, bucket))
}

func (c *Client) CountUint64(n uint64, bucket ...Field) {
	c.send(encode(MetricTypeCount, Uint64(n), c.opts.prefix, bucket))
}

func (c *Client) GaugeInt32(n int32, bucket ...Field) {
	c.send(encode(MetricTypeGauge, Int32(n), c.opts.prefix, bucket))
}

func (c *Client) GaugeUint32(n uint32, bucket ...Field) {
	c.send(encode(MetricTypeGauge, Uint32(n), c.opts.prefix, bucket))
}

func (c *Client) GaugeInt64(n int64, bucket ...Field) {
	c.send(encode(MetricTypeGauge, Int64(n), c.opts.prefix, bucket))
}

func (c *Client) GaugeUint64(n uint64, bucket ...Field) {
	c.send(encode(MetricTypeGauge, Uint64(n), c.opts.prefix, bucket))
}

func (c *Client) TimingSince(start time.Time, bucket ...Field) {
	c.send(encode(MetricTypeTiming, Float64(float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond)), c.opts.prefix, bucket))
}

func (c *Client) Timing(duration time.Duration, bucket ...Field) {
	c.send(encode(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), c.opts.prefix, bucket))
}

func (c *Client) Incrementf(template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeCount, Int32(1), c.opts.prefix, template, args))
}

func (c *Client) CountInt32f(n int32, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeCount, Int32(n), c.opts.prefix, template, args))
}

func (c *Client) CountUint32f(n uint32, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeCount, Uint32(n), c.opts.prefix, template, args))
}

func (c *Client) CountInt64f(n int64, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeCount, Int64(n), c.opts.prefix, template, args))
}

func (c *Client) CountUint64f(n uint64, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeCount, Uint64(n), c.opts.prefix, template, args))
}

func (c *Client) GaugeInt32f(n int32, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeGauge, Int32(n), c.opts.prefix, template, args))
}

func (c *Client) GaugeUint32f(n uint32, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeGauge, Uint32(n), c.opts.prefix, template, args))
}

func (c *Client) GaugeInt64f(n int64, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeGauge, Int64(n), c.opts.prefix, template, args))
}

func (c *Client) GaugeUint64f(n uint64, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeGauge, Uint64(n), c.opts.prefix, template, args))
}

func (c *Client) Timingf(duration time.Duration, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeTiming, Float64(float64(duration)/float64(time.Millisecond)), c.opts.prefix, template, args))
}

func (c *Client) TimingSincef(start time.Time, template string, args ...interface{}) {
	c.send(encodeTpl(MetricTypeTiming, Float64(float64(time.Now().Sub(start).Nanoseconds())/float64(time.Millisecond)), c.opts.prefix, template, args))
}

func (c *Client) send(b *buf) {
	if b == nil {
		return
	}
	c.cc.write(b.Bytes())
	freeBuf(b)
}
