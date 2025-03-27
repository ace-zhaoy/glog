//go:build go1.21

package slog

import (
	"context"
	"github.com/ace-zhaoy/glog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
)

type Handler struct {
	l *glog.Logger
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

	fields := make([]any, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, attr2Field(attr))
		return true
	})

	h.l.LogContext(ctx, lvl, record.Message, fields...)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	fields := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		fields = append(fields, attr2Field(attr))
	}
	cloned := h.clone()
	cloned.l = cloned.l.With(fields...)
	return cloned
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cloned := h.clone()
	cloned.l = cloned.l.With(glog.Namespace(name))
	return cloned
}

func (h *Handler) clone() *Handler {
	cloned := *h
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
		if len(attr.Value.Group()) == 0 {
			return zap.Skip()
		}
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
