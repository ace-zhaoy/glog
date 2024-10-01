package log

import (
	"context"
	"github.com/ace-zhaoy/glog"
	"sync/atomic"
	"unsafe"
)

var logger unsafe.Pointer

func Logger() *glog.Logger {
	return (*glog.Logger)(atomic.LoadPointer(&logger))
}

func SetLogger(l *glog.Logger) {
	atomic.StorePointer(&logger, unsafe.Pointer(l))
}

func init() {
	l, err := glog.NewDefault(
		glog.WithStack(glog.LevelError),
		glog.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
	SetLogger(l)
}

func WithFormatEnable() *glog.Logger {
	return Logger().WithFormatEnable()
}

func WithFormatDisable() *glog.Logger {
	return Logger().WithFormatDisable()
}

func LogContext(ctx context.Context, lvl glog.Level, msg string, args ...any) {
	Logger().LogContext(ctx, lvl, msg, args...)
}

func Log(lvl glog.Level, msg string, args ...any) {
	Logger().Log(lvl, msg, args...)
}

func Debug(msg string, args ...any) {
	Logger().Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	Logger().DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	Logger().Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger().InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	Logger().Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	Logger().WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	Logger().Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger().ErrorContext(ctx, msg, args...)
}

func Sync() error {
	return Logger().Sync()
}
