package statsd

import (
	"fmt"
	"math"
)

type MetricType uint8

const (
	MetricTypeGauge = iota
	MetricTypeCount
	MetricTypeTiming
)

type FieldType uint8

const (
	FieldTypeString = iota
	FieldTypeInt64
	FieldTypeInt32
	FieldTypeInt16
	FieldTypeInt8
	FieldTypeUint64
	FieldTypeUint32
	FieldTypeUint16
	FieldTypeUint8
	FieldTypeFloat32
	FieldTypeFloat64
)

type Field struct {
	Type FieldType
	Int  int64
	Str  string
}

func (f Field) appendTo(b *buf) {
	switch f.Type {
	case FieldTypeString:
		b.AppendString(f.Str)
	case FieldTypeFloat64:
		b.AppendFloat64(math.Float64frombits(uint64(f.Int)))
	case FieldTypeFloat32:
		b.AppendFloat32(math.Float32frombits(uint32(f.Int)))
	case FieldTypeInt8:
		b.AppendInt8(int8(f.Int))
	case FieldTypeInt16:
		b.AppendInt16(int16(f.Int))
	case FieldTypeInt32:
		b.AppendInt32(int32(f.Int))
	case FieldTypeInt64:
		b.AppendInt64(f.Int)
	case FieldTypeUint8:
		b.AppendUint8(uint8(f.Int))
	case FieldTypeUint16:
		b.AppendUint16(uint16(f.Int))
	case FieldTypeUint32:
		b.AppendUint32(uint32(f.Int))
	case FieldTypeUint64:
		b.AppendUint64(uint64(f.Int))
	default:
		panic(fmt.Sprintf("unknown field type: %v", f.Type))
	}
}

func String(val string) Field {
	return Field{Type: FieldTypeString, Str: val}
}

func Int64(val int64) Field {
	return Field{Type: FieldTypeInt64, Int: val}
}

func Int32(val int32) Field {
	return Field{Type: FieldTypeInt32, Int: int64(val)}
}

func Int16(val int16) Field {
	return Field{Type: FieldTypeInt16, Int: int64(val)}
}

func Int8(val int8) Field {
	return Field{Type: FieldTypeInt8, Int: int64(val)}
}

func Uint64(val uint64) Field {
	return Field{Type: FieldTypeUint64, Int: int64(val)}
}
func Uint32(val uint32) Field {
	return Field{Type: FieldTypeUint32, Int: int64(val)}
}
func Uint16(val uint16) Field {
	return Field{Type: FieldTypeUint16, Int: int64(val)}
}
func Uint8(val uint8) Field {
	return Field{Type: FieldTypeUint8, Int: int64(val)}
}

func Float32(val float32) Field {
	return Field{Type: FieldTypeFloat32, Int: int64(math.Float32bits(val))}
}

func Float64(val float64) Field {
	return Field{Type: FieldTypeFloat64, Int: int64(math.Float64bits(val))}
}

func encode(typ MetricType, val Field, prefix string, bucket []Field) *buf {
	n := len(bucket)
	if n == 0 {
		return nil
	}

	b := getBuf()

	if prefix != "" {
		b.AppendString(prefix)
		b.AppendString(".")
	}

	last := n - 1
	for i := range bucket {
		bucket[i].appendTo(b)
		if i < last {
			b.AppendString(".")
		}
	}

	b.AppendString(":")
	val.appendTo(b)

	b.AppendString("|")
	switch typ {
	case MetricTypeGauge:
		b.AppendString("g\n")
	case MetricTypeCount:
		b.AppendString("c\n")
	case MetricTypeTiming:
		b.AppendString("ms\n")
	default:
		panic(fmt.Sprintf("unknown field type: %v", typ))
	}

	return b
}

func encodeTpl(typ MetricType, val Field, prefix string, template string, fmtArgs []interface{}) *buf {
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	return encode(typ, val, prefix, []Field{String(msg)})
}
