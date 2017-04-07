package statsd_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/kirk91/statsd"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	_, err := statsd.New("foo", "127.0.0.1:1")
	assert.Error(t, err)
	_, err = statsd.New("udp", "127.0.0.1:1")
	assert.NoError(t, err)

	l, err := net.ListenPacket("udp", "")
	assert.NoError(t, err)
	assert.NotNil(t, l)
	c, err := statsd.New("udp", l.LocalAddr().String())
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

type mockServer struct {
	l      net.PacketConn
	closed bool
	buf    bytes.Buffer
}

func newMockServer(t *testing.T) *mockServer {
	l, err := net.ListenPacket("udp", "")
	assert.NoError(t, err)

	s := &mockServer{l: l}

	go func() {
		b := make([]byte, 1024)
		for {
			if s.closed {
				return
			}
			n, _, err := l.ReadFrom(b)
			if n == 0 {
				continue
			}
			assert.NoError(t, err)
			s.buf.Write(b[:n])
		}
	}()

	return s
}

func (s *mockServer) Addr() string {
	return s.l.LocalAddr().String()
}

func (s *mockServer) Close() {
	s.closed = true
	s.l.Close()
}

func (s *mockServer) Reset() {
	s.buf.Reset()
}

func (s *mockServer) Content() string {
	return string(s.buf.Bytes())
}

func TestIncrement(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.FlushPeriod(time.Nanosecond*500))

	c.Increment(statsd.Int8(1), statsd.Int16(200), statsd.Int32(1000), statsd.Int64(10000))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "1.200.1000.10000:1|c\n", s.Content())

	s.Reset()
	c.Increment(statsd.Uint8(1), statsd.Uint16(200), statsd.Uint32(1000), statsd.Uint64(10000))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "1.200.1000.10000:1|c\n", s.Content())

	s.Reset()
	c.Increment(statsd.String("foo"))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "foo:1|c\n", s.Content())

	s.Reset()
	c.Increment(statsd.String("foo"), statsd.String("bar"))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "foo.bar:1|c\n", s.Content())
}

func TestCount(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.FlushPeriod(time.Nanosecond*500))

	c.CountInt32(1, statsd.String("foo"))
	c.CountInt64(3, statsd.String("foo"))
	c.CountInt32(10, statsd.String("bar"))
	c.CountInt64(100, statsd.String("bar"))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "foo:1|c\nfoo:3|c\nbar:10|c\nbar:100|c\n", s.Content())
}

func TestGauge(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.FlushPeriod(time.Nanosecond*500))

	c.GaugeInt32(1)
	c.GaugeInt32(1, statsd.String("foo"), statsd.String("bar"))
	c.GaugeInt64(2, statsd.String("foo"), statsd.String("bar"))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "foo.bar:1|g\nfoo.bar:2|g\n", s.Content())
}

func TestTiming(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.FlushPeriod(time.Nanosecond*500))
	c.Timing(time.Now(), statsd.String("foo"))
	time.Sleep(time.Millisecond * 100)
	assert.Contains(t, s.Content(), "|ms\n")
}

func TestPrefix(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.Prefix("juju"))
	c.Increment(statsd.String("foo"))
	c.Increment(statsd.String("bar"))
	c.CountInt32(3, statsd.String("zoo"))
	c.CountInt64(10, statsd.String("kong"))
	c.GaugeInt32(100, statsd.String("mong"))
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, "juju.foo:1|c\njuju.bar:1|c\njuju.zoo:3|c\njuju.kong:10|c\njuju.mong:100|g\n", s.Content())
}

func TestMaxPacketSize(t *testing.T) {
	s := newMockServer(t)
	defer s.Close()

	c, _ := statsd.New("udp", s.Addr(), statsd.MaxPacketSize(20))
	c.Increment(statsd.String("foo.bar.zoo"))
	c.Increment(statsd.String("foo.bar.zoo"))
	time.Sleep(time.Millisecond * 80)
	assert.Equal(t, "foo.bar.zoo:1|c\n", s.Content())
}

func TestErrorHandler(t *testing.T) {
	l, err := net.Listen("tcp", "")
	assert.NoError(t, err)

	var gotErr bool
	c, _ := statsd.New("tcp", l.Addr().String(), statsd.ErrorHandler(func(error) {
		gotErr = true
	}), statsd.FlushPeriod(50*time.Nanosecond))
	c.Increment(statsd.String("foo.bar.zoo"))
	l.Close()
	time.Sleep(time.Millisecond * 100)
	assert.True(t, gotErr)
}

func BenchmarkIncrement(b *testing.B) {
	c, _ := statsd.New("udp", "127.0.0.1:1")
	foo, bar, zoo := "foo", "bar", 1

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Increment(statsd.String(foo), statsd.String(bar), statsd.Int32(int32(zoo)))
	}
}

func BenchmarkCount(b *testing.B) {
	c, _ := statsd.New("udp", "127.0.0.1:1")
	foo, bar, zoo := "foo", "bar", int32(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.CountInt32(10, statsd.String(foo), statsd.String(bar), statsd.Int32(zoo))
	}
}

func BenchmarkGauge(b *testing.B) {
	c, _ := statsd.New("udp", "127.0.0.1:1")
	foo, bar, zoo := "foo", "bar", int64(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.GaugeInt32(10, statsd.String(foo), statsd.String(bar), statsd.Int64(zoo))
	}
}

func BenchmarkTiming(b *testing.B) {
	c, _ := statsd.New("udp", "127.0.0.1:1")
	foo, bar, zoo := "foo", "bar", int32(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Timing(time.Now(), statsd.String(foo), statsd.String(bar), statsd.Int32(zoo))
	}
}