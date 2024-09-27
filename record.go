package glog

import "context"

// NewRecordWithCapacity returns a new Record with the given capacity.
func NewRecordWithCapacity(capacity int) *Record {
	return &Record{fields: make([]Field, 0, capacity)}
}

type Record struct {
	fields []Field
}

func (r *Record) AddFields(fields ...Field) {
	r.fields = append(r.fields, fields...)
}

func (r *Record) Add(args ...any) {
	r.fields = append(r.fields, argsToFields(args)...)
}

func (r *Record) Fields() []Field {
	return r.fields
}

type ContextHandler func(ctx context.Context, record *Record)

// BuildContextHandler builds a ContextHandler from key and alias.
// key is the context key for the value.
// alias is an optional alias for the field key.
func BuildContextHandler(key string, alias ...string) ContextHandler {
	fieldName := key
	if len(alias) > 0 && alias[0] != "" {
		fieldName = alias[0]
	}
	return func(ctx context.Context, record *Record) {
		if v := ctx.Value(key); v != nil {
			record.AddFields(Any(fieldName, v))
		}
	}
}

const (
	badKey  = "!BADKEY"
	noValue = "!NOVALUE"
)

// argsToFields converts arguments to fields.
func argsToFields(args []any) (fields []Field) {
	fields = make([]Field, 0, len(args))
	if len(args) == 0 {
		return
	}
	argsLen := len(args)
	for i := 0; i < argsLen; i++ {
		switch v := args[i].(type) {
		case Field:
			fields = append(fields, v)
		case string:
			if i == argsLen-1 {
				fields = append(fields, String(noValue, v))
				return
			}
			fields = append(fields, Any(v, args[i+1]))
			i += 1
		default:
			fields = append(fields, Any(badKey, v))
		}
	}
	return
}
