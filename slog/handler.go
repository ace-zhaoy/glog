//go:build go1.21

package slog

import (
	"context"
	"github.com/ace-zhaoy/glog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"slices"
)

type Handler struct {
	l *glog.Logger

	groups []string
}

var _ slog.Handler = (*Handler)(nil)

func NewHandler(l *glog.Logger) *Handler {
	return &Handler{
		l: l,
	}
}

func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	return h.l.Enabled(LevelConverter(lvl))
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	lvl := LevelConverter(record.Level)
	if !h.l.Enabled(lvl) {
		return nil
	}
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	fields := h.toFields(attrs)
	h.l.LogContext(ctx, lvl, record.Message, fields...)
	return nil
}

func isEmptyGroup(attr slog.Attr) bool {
	if attr.Value.Kind() != slog.KindGroup {
		return false
	}

	return len(attr.Value.Group()) == 0
}

func (h *Handler) toFields(attrs []slog.Attr) []any {
	if len(attrs) == 0 {
		return nil
	}
	fields, index := make([]any, len(attrs)+len(h.groups)), len(h.groups)
	for _, v := range attrs {
		if isEmptyGroup(v) {
			continue
		}
		field := attr2Field(v)
		if field.Equals(glog.Skip()) {
			continue
		}
		fields[index] = field
		index++
	}

	if index == len(h.groups) {
		return nil
	}

	for i, v := range h.groups {
		fields[i] = glog.Namespace(v)
	}
	return fields[:index]
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := h.toFields(attrs)
	if len(fields) == 0 {
		return h
	}

	cloned := h.clone()
	cloned.l = cloned.l.With(fields...)
	cloned.groups = nil
	return cloned
}

func (h *Handler) WithGroup(name string) slog.Handler {
	cloned := h.clone()
	cloned.groups = append(cloned.groups, name)
	return cloned
}

func (h *Handler) clone() *Handler {
	cloned := *h
	cloned.groups = slices.Clip(h.groups)
	return &cloned
}

func LevelConverter(lvl slog.Level) glog.Level {
	switch {
	case lvl >= slog.LevelError:
		return glog.LevelError
	case lvl >= slog.LevelWarn:
		return glog.LevelWarn
	case lvl >= slog.LevelInfo:
		return glog.LevelInfo
	default:
		return glog.LevelDebug
	}
}

func attr2Field(attr slog.Attr) glog.Field {
	if attr.Equal(slog.Attr{}) {
		return zap.Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindGroup:
		val := group(attr.Value.Group())
		if attr.Key == "" {
			return zap.Inline(val)
		}
		return zap.Object(attr.Key, val)
	case slog.KindLogValuer:
		return zap.Inline(logValuer{attr})
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

type logValuer struct {
	attr slog.Attr
}

func (lv logValuer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	attr2Field(slog.Attr{
		Key:   lv.attr.Key,
		Value: lv.attr.Value.Resolve(),
	}).AddTo(enc)
	return nil
}

type group []slog.Attr

func (g group) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range g {
		attr2Field(attr).AddTo(enc)
	}
	return nil
}
