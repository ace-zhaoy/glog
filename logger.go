package glog

import (
	"context"
	"fmt"
	"github.com/ace-zhaoy/glog/stacktrace"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	callerSkipOffset = 3
)

type Logger struct {
	core       Core
	name       string
	addCaller  bool
	stackLevel LevelEnabler
	callerSkip int

	formatEnabled   bool
	contextHandlers []ContextHandler
}

func NewLogger(core Core, opts ...Option) *Logger {
	l := &Logger{
		core: core,
	}
	return l.WithOptions(opts...)
}

func (l *Logger) clone() *Logger {
	c := *l
	contextHandlers := make([]ContextHandler, len(l.contextHandlers))
	copy(contextHandlers, l.contextHandlers)
	c.contextHandlers = contextHandlers
	return &c
}

func (l *Logger) formatMessage(msg string, args []any) (message string, formatted bool) {
	message = msg
	if !l.formatEnabled || len(args) == 0 {
		return
	}

	for i := 0; i < len(args); i++ {
		if _, ok := args[i].(Field); ok {
			return
		}
	}

	if countPercent(msg) != len(args) {
		return
	}

	return fmt.Sprintf(msg, args...), true
}

func (l *Logger) With(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}
	log := l.clone()
	log.core = l.core.With(argsToFields(args))
	return log
}

func (l *Logger) WithOptions(opts ...Option) *Logger {
	if len(opts) == 0 {
		return l
	}
	log := l.clone()
	for _, opt := range opts {
		opt.apply(log)
	}
	return log
}

func (l *Logger) WithFormat(formatEnabled bool) *Logger {
	if l.formatEnabled == formatEnabled {
		return l
	}
	log := l.clone()
	log.formatEnabled = formatEnabled
	return log
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil || len(l.contextHandlers) == 0 {
		return l
	}

	record := NewRecordWithCapacity(len(l.contextHandlers))
	for _, handler := range l.contextHandlers {
		handler(ctx, record)
	}

	log := l.clone()
	log.core = l.core.With(record.Fields())

	return log
}

func (l *Logger) check(lvl Level, msg string) (ce *zapcore.CheckedEntry) {
	ent := zapcore.Entry{
		LoggerName: l.name,
		Time:       time.Now(),
		Level:      lvl,
		Message:    msg,
	}
	ce = l.core.Check(ent, nil)
	if ce == nil {
		return
	}

	addStack := false
	if l.stackLevel != nil {
		addStack = l.stackLevel.Enabled(ce.Level)
	}
	if !l.addCaller && !addStack {
		return
	}

	stackDepth := stacktrace.First
	if addStack {
		stackDepth = stacktrace.Full
	}

	stack := stacktrace.Capture(l.callerSkip+callerSkipOffset, stackDepth)
	defer stack.Free()

	frame, more := stack.Next()
	if l.addCaller {
		ce.Caller = zapcore.EntryCaller{
			Defined:  frame.PC != 0,
			PC:       frame.PC,
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		}
	}

	if addStack {
		formatter := stacktrace.GetFormatter()
		defer formatter.Free()

		formatter.FormatFrame(frame)
		if more {
			formatter.FormatStack(stack)
		}

		ce.Stack = formatter.String()
	}
	return ce
}

func (l *Logger) log(ctx context.Context, lvl Level, msg string, args ...any) {
	if !l.core.Enabled(lvl) {
		return
	}

	msg, msgFormatted := l.formatMessage(msg, args)
	ce := l.check(lvl, msg)
	if ce == nil {
		return
	}

	var fields []Field
	capacity := len(l.contextHandlers)
	if !msgFormatted {
		fields = argsToFields(args)
		capacity += len(fields)
	}

	record := NewRecordWithCapacity(capacity)
	if ctx != nil && len(l.contextHandlers) > 0 {
		for _, handler := range l.contextHandlers {
			handler(ctx, record)
		}
	}
	record.AddFields(fields...)

	ce.Write(record.Fields()...)
}

func (l *Logger) LogContext(ctx context.Context, lvl Level, msg string, args ...any) {
	l.log(ctx, lvl, msg, args...)
}

func (l *Logger) Log(lvl Level, msg string, args ...any) {
	l.log(nil, lvl, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(nil, LevelDebug, msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(nil, LevelInfo, msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(nil, LevelWarn, msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(nil, LevelError, msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelError, msg, args...)
}

func (l *Logger) Sync() error {
	return l.core.Sync()
}

func countPercent(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			if i+1 < len(s) && s[i+1] == '%' {
				i++
			} else {
				count++
			}
		}
	}
	return count
}
