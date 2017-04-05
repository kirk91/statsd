package statsd

import (
	"strconv"
	"sync"
)

type buf struct {
	bs []byte
}

func (b *buf) Bytes() []byte {
	return b.bs
}

func (b *buf) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

func (b *buf) AppendUint(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10)
}

func (b *buf) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}

func (b *buf) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &buf{bs: make([]byte, 0, 512)}
	},
}

func getBuf() *buf {
	return bufPool.Get().(*buf)
}

func freeBuf(b *buf) {
	b.bs = b.bs[:0]
	bufPool.Put(b)
}
