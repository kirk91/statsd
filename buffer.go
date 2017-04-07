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

func (b *buf) AppendInt64(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}

func (b *buf) AppendUint64(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10)
}

func (b *buf) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}

func (b *buf) AppendString(s string) {
	b.bs = append(b.bs, s...)
}

func (b *buf) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
}

func (b *buf) AppendFloat32(v float32) { b.AppendFloat(float64(v), 32) }
func (b *buf) AppendFloat64(v float64) { b.AppendFloat(v, 64) }
func (b *buf) AppendInt8(v int8)       { b.AppendInt64(int64(v)) }
func (b *buf) AppendInt16(v int16)     { b.AppendInt64(int64(v)) }
func (b *buf) AppendInt32(v int32)     { b.AppendInt64(int64(v)) }
func (b *buf) AppendUint8(v uint8)     { b.AppendUint64(uint64(v)) }
func (b *buf) AppendUint16(v uint16)   { b.AppendUint64(uint64(v)) }
func (b *buf) AppendUint32(v uint32)   { b.AppendUint64(uint64(v)) }

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
