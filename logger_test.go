package glog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogger(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	assert.NotNil(t, logger, "Expected new logger to be created")
	assert.Equal(t, core, logger.core, "Expected core to be set")
}

func TestLogger_clone(t *testing.T) {
	h1 := func(ctx context.Context, record *Record) {}
	l1 := &Logger{contextHandlers: []ContextHandler{h1}}
	l2 := l1.clone()

	assert.NotEqual(t, l1.contextHandlers, l2.contextHandlers, "Expected contextHandlers to be different")
	assert.NotEqual(t, l1, l2, "Expected l1 and l2 to be different")
	assert.Equal(t, len(l2.contextHandlers), 1, "Expected contextHandlers to be cloned")
}

func TestLogger_formatMessage(t *testing.T) {
	type myLogger struct {
		core            Core
		name            string
		addCaller       bool
		stackLevel      LevelEnabler
		callerSkip      int
		formatEnabled   bool
		contextHandlers []ContextHandler
	}
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name          string
		logger        myLogger
		args          args
		wantMessage   string
		wantFormatted bool
	}{
		{
			name:   "formatEnabled is false, no args",
			logger: myLogger{},
			args: args{
				msg: "msg1",
			},
			wantMessage:   "msg1",
			wantFormatted: false,
		},
		{
			name:   "formatEnabled is false, with args",
			logger: myLogger{},
			args: args{
				msg:  "msg is %s",
				args: []any{"test"},
			},
			wantMessage:   "msg is %s",
			wantFormatted: false,
		},
		{
			name:   "formatEnabled is true, no args",
			logger: myLogger{formatEnabled: true},
			args: args{
				msg: "msg1",
			},
			wantMessage:   "msg1",
			wantFormatted: false,
		},
		{
			name:   "formatEnabled is true, with args",
			logger: myLogger{formatEnabled: true},
			args: args{
				msg:  "msg is %s",
				args: []any{"test"},
			},
			wantMessage:   "msg is test",
			wantFormatted: true,
		},
		{
			name:   "formatEnabled is true, with multiple args",
			logger: myLogger{formatEnabled: true},
			args: args{
				msg:  "msg is %s",
				args: []any{"test", "test2"},
			},
			wantMessage:   "msg is %s",
			wantFormatted: false,
		},
		{
			name:   "formatEnabled is true, with Filed args",
			logger: myLogger{formatEnabled: true},
			args: args{
				msg:  "msg is %s",
				args: []any{String("k", "v")},
			},
			wantMessage:   "msg is %s",
			wantFormatted: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				core:            tt.logger.core,
				name:            tt.logger.name,
				addCaller:       tt.logger.addCaller,
				stackLevel:      tt.logger.stackLevel,
				callerSkip:      tt.logger.callerSkip,
				formatEnabled:   tt.logger.formatEnabled,
				contextHandlers: tt.logger.contextHandlers,
			}
			gotMessage, gotFormatted := l.formatMessage(tt.args.msg, tt.args.args)
			assert.Equalf(t, tt.wantMessage, gotMessage, "formatMessage(%v, %v)", tt.args.msg, tt.args.args)
			assert.Equalf(t, tt.wantFormatted, gotFormatted, "formatMessage(%v, %v)", tt.args.msg, tt.args.args)
		})
	}
}

func TestLogger_With(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	logWithFields := logger.With("key1", "value1", "key2", "value2", String("key3", "value3"))
	assert.NotNil(t, logWithFields, "Expected logger with fields")
	assert.Equal(t, 3, len(core.fields), "Expected fields to be added to core")
}

func TestLogger_WithOptions(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	opt := WithName("test_logger")
	loggerWithOptions := logger.WithOptions(opt)
	assert.Equal(t, "test_logger", loggerWithOptions.name, "Expected logger name to be set")
}

func TestLogger_WithFormat(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	loggerWithFormat := logger.WithFormat(true)
	assert.True(t, loggerWithFormat.formatEnabled, "Expected format to be enabled")

	loggerWithFormat = logger.WithFormat(false)
	assert.False(t, loggerWithFormat.formatEnabled, "Expected format to be disabled")
}

func TestLogger_WithFormatEnable(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	loggerWithFormat := logger.WithFormatEnable()
	assert.True(t, loggerWithFormat.formatEnabled, "Expected format to be enabled")
}

func TestLogger_WithFormatDisable(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	loggerWithFormat := logger.WithFormatDisable()
	assert.False(t, loggerWithFormat.formatEnabled, "Expected format to be disabled")
}

func TestLogger_WithContext(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(
		core,
		WithContextHandlers(
			BuildContextHandler("k1"),
			BuildContextHandler("key2", "k2"),
			func(ctx context.Context, record *Record) {
				key := "k3"
				if v := ctx.Value(key); v != nil {
					record.AddFields(Any(key, v))
				}
			},
		),
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "k1", "v1")
	ctx = context.WithValue(ctx, "key2", 2)
	logger = logger.WithContext(ctx)

	assert.Equal(t, core.fields, []Field{
		String("k1", "v1"),
		Int("k2", 2),
	}, "Expected context to be added to core")
}

func TestLogger_Enabled(t *testing.T) {
	core := &mockCore{enabled: true}
	logger := NewLogger(core)

	assert.True(t, logger.Enabled(LevelInfo), "Expected logger to be enabled")
}

func Test_countPercent(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "No format here",
			args: args{
				s: "No format here",
			},
			want: 0,
		},
		{
			name: "Hello %s, how are you?",
			args: args{
				s: "Hello %s, how are you?",
			},
			want: 1,
		},
		{
			name: "Hello %%s, how are you?",
			args: args{
				s: "Hello %%s, how are you?",
			},
			want: 0,
		},
		{
			name: "Hello %s%%, how are you?",
			args: args{
				s: "Hello %s%%, how are you?",
			},
			want: 1,
		},
		{
			name: "Hello %%s%%, how are you?",
			args: args{
				s: "Hello %%s%%, how are you?",
			},
			want: 0,
		},
		{
			name: "Hello %%%s, how are you?",
			args: args{
				s: "Hello %%%s, how are you?",
			},
			want: 1,
		},
		{
			name: "Hello %%%%s, how are you?",
			args: args{
				s: "Hello %%%%s, how are you?",
			},
			want: 0,
		},
		{
			name: "Empty string",
			args: args{
				s: "",
			},
			want: 0,
		},
		{
			name: "Only percent signs",
			args: args{
				s: "%%%%",
			},
			want: 0,
		},
		{
			name: "Single percent sign",
			args: args{
				s: "%",
			},
			want: 1,
		},
		{
			name: "Multiple percent signs",
			args: args{
				s: "%%%s%%%d",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, countPercent(tt.args.s), "countPercent(%v)", tt.args.s)
		})
	}
}

func TestLogger_check(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core)

	t.Run("check returns nil if core is not enabled", func(t *testing.T) {
		core.enabled = false
		ce := logger.check(LevelInfo, "test message")
		assert.Nil(t, ce, "Expected check to return nil when core is not enabled")
	})

	t.Run("check returns CheckedEntry if core is enabled", func(t *testing.T) {
		core.enabled = true
		ce := logger.check(LevelInfo, "test message")
		assert.NotNil(t, ce, "Expected check to return CheckedEntry when core is enabled")
		assert.Equal(t, "test message", ce.Message, "Expected message to be set in CheckedEntry")
	})

	t.Run("check adds caller information if addCaller is true", func(t *testing.T) {
		logger.addCaller = true
		ce := logger.check(LevelInfo, "test message")
		assert.NotNil(t, ce.Caller, "Expected caller information to be added")
		assert.True(t, ce.Caller.Defined, "Expected caller to be defined")
	})

	t.Run("check adds stack trace if stackLevel is enabled", func(t *testing.T) {
		logger.stackLevel = LevelEnablerFunc(func(lvl Level) bool {
			return lvl == LevelInfo
		})
		ce := logger.check(LevelInfo, "test message")
		assert.NotEmpty(t, ce.Stack, "Expected stack trace to be added")
	})
}

func TestLogger_log(t *testing.T) {
	core := &mockCore{}
	logger := NewLogger(core, WithContextHandlers(BuildContextHandler("key")))

	t.Run("log without context", func(t *testing.T) {
		core.reset()
		core.enabled = true
		logger.log(nil, LevelInfo, "test message", String("key", "value"))
		assert.Len(t, core.entries, 1, "Expected one log entry")
		assert.Equal(t, "test message", core.entries[0].Message, "Expected message to be 'test message'")
		assert.Contains(t, core.fields, String("key", "value"), "Expected fields to contain 'key: value'")
	})

	t.Run("log with context", func(t *testing.T) {
		core.reset()
		core.enabled = true
		ctx := context.WithValue(context.Background(), "key", "value")
		logger.log(ctx, LevelInfo, "test message")
		assert.Len(t, core.entries, 1, "Expected one log entry")
		assert.Equal(t, "test message", core.entries[0].Message, "Expected message to be 'test message'")
		assert.Contains(t, core.fields, String("key", "value"), "Expected fields to contain 'key: value'")
	})

	t.Run("log with formatted message", func(t *testing.T) {
		core.reset()
		core.enabled = true
		logger = logger.WithFormatEnable()
		logger.log(nil, LevelInfo, "Hello %s", "world")
		assert.Len(t, core.entries, 1, "Expected one log entry")
		assert.Equal(t, "Hello world", core.entries[0].Message, "Expected formatted message to be 'Hello world'")
	})

	t.Run("log with stack trace", func(t *testing.T) {
		core.reset()
		core.enabled = true
		logger.stackLevel = LevelEnablerFunc(func(lvl Level) bool {
			return lvl == LevelInfo
		})
		logger.log(nil, LevelInfo, "test message")
		assert.Len(t, core.entries, 1, "Expected one log entry")
		assert.NotEmpty(t, core.entries[0].Stack, "Expected stack trace to be added")
	})

	t.Run("log with caller information", func(t *testing.T) {
		core.reset()
		core.enabled = true
		logger.addCaller = true
		logger.log(nil, LevelInfo, "test message")
		assert.Len(t, core.entries, 1, "Expected one log entry")
		assert.True(t, core.entries[0].Caller.Defined, "Expected caller information to be added")
	})
}
