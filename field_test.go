package glog

import (
	"github.com/ace-zhaoy/glog/stacktrace"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regexp"
	"testing"
	"time"
)

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

func TestFields(t *testing.T) {
	var (
		boolVal       = bool(true)
		complex128Val = complex128(complex(0, 0))
		complex64Val  = complex64(complex(0, 0))
		durationVal   = time.Duration(time.Second)
		float64Val    = float64(1.0)
		float32Val    = float32(1.0)
		intVal        = int(1)
		int64Val      = int64(1)
		int32Val      = int32(1)
		int16Val      = int16(1)
		int8Val       = int8(1)
		stringVal     = string("hello")
		timeVal       = time.Unix(100000, 0)
		uintVal       = uint(1)
		uint64Val     = uint64(1)
		uint32Val     = uint32(1)
		uint16Val     = uint16(1)
		uint8Val      = uint8(1)
		uintptrVal    = uintptr(1)
		name          = username("hello")
	)

	tests := []struct {
		name  string
		field Field
		want  zapcore.Field
	}{
		{
			name:  "Skip",
			field: Skip(),
			want:  zap.Skip(),
		},
		{
			name:  "Binary",
			field: Binary("k", []byte("Binary")),
			want:  zap.Binary("k", []byte("Binary")),
		},
		{
			name:  "Bool",
			field: Bool("k", boolVal),
			want:  zap.Bool("k", boolVal),
		},
		{
			name:  "Boolp",
			field: Boolp("k", &boolVal),
			want:  zap.Boolp("k", &boolVal),
		},
		{
			name:  "ByteString",
			field: ByteString("k", []byte("ByteString")),
			want:  zap.ByteString("k", []byte("ByteString")),
		},
		{
			name:  "Complex128",
			field: Complex128("k", complex128Val),
			want:  zap.Complex128("k", complex128Val),
		},
		{
			name:  "Complex128p",
			field: Complex128p("k", &complex128Val),
			want:  zap.Complex128p("k", &complex128Val),
		},
		{
			name:  "Complex64",
			field: Complex64("k", complex64Val),
			want:  zap.Complex64("k", complex64Val),
		},
		{
			name:  "Complex64p",
			field: Complex64p("k", &complex64Val),
			want:  zap.Complex64p("k", &complex64Val),
		},
		{
			name:  "Float64",
			field: Float64("k", float64Val),
			want:  zap.Float64("k", float64Val),
		},
		{
			name:  "Float64p",
			field: Float64p("k", &float64Val),
			want:  zap.Float64p("k", &float64Val),
		},
		{
			name:  "Float32",
			field: Float32("k", float32Val),
			want:  zap.Float32("k", float32Val),
		},
		{
			name:  "Float32p",
			field: Float32p("k", &float32Val),
			want:  zap.Float32p("k", &float32Val),
		},
		{
			name:  "Int",
			field: Int("k", intVal),
			want:  zap.Int("k", intVal),
		},
		{
			name:  "Intp",
			field: Intp("k", &intVal),
			want:  zap.Intp("k", &intVal),
		},
		{
			name:  "Int64",
			field: Int64("k", int64Val),
			want:  zap.Int64("k", int64Val),
		},
		{
			name:  "Int64p",
			field: Int64p("k", &int64Val),
			want:  zap.Int64p("k", &int64Val),
		},
		{
			name:  "Int32",
			field: Int32("k", int32Val),
			want:  zap.Int32("k", int32Val),
		},
		{
			name:  "Int32p",
			field: Int32p("k", &int32Val),
			want:  zap.Int32p("k", &int32Val),
		},
		{
			name:  "Int16",
			field: Int16("k", int16Val),
			want:  zap.Int16("k", int16Val),
		},
		{
			name:  "Int16p",
			field: Int16p("k", &int16Val),
			want:  zap.Int16p("k", &int16Val),
		},
		{
			name:  "Int8",
			field: Int8("k", int8Val),
			want:  zap.Int8("k", int8Val),
		},
		{
			name:  "Int8p",
			field: Int8p("k", &int8Val),
			want:  zap.Int8p("k", &int8Val),
		},
		{
			name:  "String",
			field: String("k", stringVal),
			want:  zap.String("k", stringVal),
		},
		{
			name:  "Stringp",
			field: Stringp("k", &stringVal),
			want:  zap.Stringp("k", &stringVal),
		},
		{
			name:  "Uint",
			field: Uint("k", uintVal),
			want:  zap.Uint("k", uintVal),
		},
		{
			name:  "Uintp",
			field: Uintp("k", &uintVal),
			want:  zap.Uintp("k", &uintVal),
		},
		{
			name:  "Uint64",
			field: Uint64("k", uint64Val),
			want:  zap.Uint64("k", uint64Val),
		},
		{
			name:  "Uint64p",
			field: Uint64p("k", &uint64Val),
			want:  zap.Uint64p("k", &uint64Val),
		},
		{
			name:  "Uint32",
			field: Uint32("k", uint32Val),
			want:  zap.Uint32("k", uint32Val),
		},
		{
			name:  "Uint32p",
			field: Uint32p("k", &uint32Val),
			want:  zap.Uint32p("k", &uint32Val),
		},
		{
			name:  "Uint16",
			field: Uint16("k", uint16Val),
			want:  zap.Uint16("k", uint16Val),
		},
		{
			name:  "Uint16p",
			field: Uint16p("k", &uint16Val),
			want:  zap.Uint16p("k", &uint16Val),
		},
		{
			name:  "Uint8",
			field: Uint8("k", uint8Val),
			want:  zap.Uint8("k", uint8Val),
		},
		{
			name:  "Uint8p",
			field: Uint8p("k", &uint8Val),
			want:  zap.Uint8p("k", &uint8Val),
		},
		{
			name:  "Uintptr",
			field: Uintptr("k", uintptrVal),
			want:  zap.Uintptr("k", uintptrVal),
		},
		{
			name:  "Uintptrp",
			field: Uintptrp("k", &uintptrVal),
			want:  zap.Uintptrp("k", &uintptrVal),
		},
		{
			name:  "Reflect",
			field: Reflect("k", nil),
			want:  zap.Reflect("k", nil),
		},
		{
			name:  "Namespace",
			field: Namespace("k"),
			want:  zap.Namespace("k"),
		},
		{
			name:  "Stringer",
			field: Stringer("k", nil),
			want:  zap.Stringer("k", nil),
		},
		{
			name:  "Time",
			field: Time("k", timeVal),
			want:  zap.Time("k", timeVal),
		},
		{
			name:  "Timep",
			field: Timep("k", &timeVal),
			want:  zap.Timep("k", &timeVal),
		},
		{
			name:  "Duration",
			field: Duration("k", durationVal),
			want:  zap.Duration("k", durationVal),
		},
		{
			name:  "Durationp",
			field: Durationp("k", &durationVal),
			want:  zap.Durationp("k", &durationVal),
		},
		{
			name:  "Object",
			field: Object("k", name),
			want:  zap.Object("k", name),
		},
		{
			name:  "Inline",
			field: Inline(name),
			want:  zap.Inline(name),
		},
		{
			name:  "Any",
			field: Any("k", "any"),
			want:  zap.Any("k", "any"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.field, "want: %+#v, got: %+#v", tt.want, tt.field)
		})
	}
}

func TestStackField(t *testing.T) {
	f := Stack("stacktrace")
	assert.Equal(t, "stacktrace", f.Key, "Unexpected field key.")
	assert.Equal(t, zapcore.StringType, f.Type, "Unexpected field type.")
	r := regexp.MustCompile(`field_test.go:(\d+)`)
	assert.Equal(t, r.ReplaceAllString(stacktrace.Take(0), "field_test.go"), r.ReplaceAllString(f.String, "field_test.go"), "Unexpected stack trace")
}

func TestStackSkipField(t *testing.T) {
	f := StackSkip("stacktrace", 0)
	assert.Equal(t, "stacktrace", f.Key, "Unexpected field key.")
	assert.Equal(t, zapcore.StringType, f.Type, "Unexpected field type.")
	r := regexp.MustCompile(`field_test.go:(\d+)`)
	assert.Equal(t, r.ReplaceAllString(stacktrace.Take(0), "field_test.go"), r.ReplaceAllString(f.String, "field_test.go"), f.String, "Unexpected stack trace")
}
