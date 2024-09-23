package cores

import (
	"go.uber.org/zap/zapcore"
)

type LazyCore struct {
	core zapcore.Core

	fields []zapcore.Field
}

var _ zapcore.Core = (*LazyCore)(nil)

func NewLazyCore(core zapcore.Core, fields ...zapcore.Field) *LazyCore {
	return &LazyCore{
		core:   core,
		fields: fields,
	}
}

func (l *LazyCore) clone() *LazyCore {
	c := *l
	fields := make([]zapcore.Field, len(l.fields))
	copy(fields, l.fields)
	c.fields = fields
	return &c
}

func (l *LazyCore) Enabled(lvl zapcore.Level) bool {
	return l.core.Enabled(lvl)
}

func (l *LazyCore) With(fields []zapcore.Field) zapcore.Core {
	c := l.clone()
	c.fields = append(c.fields, fields...)
	return c
}

func (l *LazyCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return l.core.Check(ent, ce)
}

func (l *LazyCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return l.core.Write(ent, append(l.fields, fields...))
}

func (l *LazyCore) Sync() error {
	return l.core.Sync()
}
