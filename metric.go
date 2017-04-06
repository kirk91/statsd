package statsd

import "fmt"

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
	FieldTypeInt
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
	case FieldTypeInt64, FieldTypeInt32, FieldTypeInt16, FieldTypeInt8, FieldTypeInt:
		b.AppendInt(f.Int)
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

func Int(val int) Field {
	return Field{Type: FieldTypeInt, Int: int64(val)}
}

func encode(typ MetricType, prefix string, bucket []Field, val Field) *buf {
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
