package statsd

import (
	"net"
	"sync"
	"time"
)

type clientConn struct {
	network, addr string
	c             *Client
	conn          net.Conn

	mu  sync.Mutex
	buf []byte
}

func newClientConn(network, addr string, c *Client) (*clientConn, error) {
	conn, err := net.DialTimeout(network, addr, c.opts.timeout)
	if err != nil {
		return nil, err
	}

	// When using UDP do a quick check to see if something is listening on the
	// given port to return an error as soon as possible.
	if network == "udp" {
		for i := 0; i < 2; i++ {
			_, err = conn.Write(nil)
			if err != nil {
				conn.Close()
				return nil, err
			}
		}
	}

	cc := &clientConn{
		network: network,
		addr:    addr,
		c:       c,
		buf:     make([]byte, 0, c.opts.maxPacketSize),
	}

	go func() {
		ticker := time.NewTicker(c.opts.flushPeriod)
		for _ = range ticker.C {
			cc.mu.Lock()
			cc.flush()
			cc.mu.Unlock()
		}
	}()

	return cc, nil
}

func (cc *clientConn) write(b []byte) {
	cc.mu.Lock()
	if len(cc.buf)+len(b) > cap(cc.buf) {
		cc.flush()
	}
	cc.buf = append(cc.buf, b...)
	cc.mu.Unlock()
}

func (cc *clientConn) flush() {
	if len(cc.buf) == 0 {
		return
	}
	_, err := cc.conn.Write(cc.buf[:len(cc.buf)])
	cc.handleError(err)
	cc.buf = cc.buf[:0]
}

func (cc *clientConn) handleError(err error) {
	if err != nil && cc.c.opts.errHandler != nil {
		cc.c.opts.errHandler(err)
	}
}
