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
	core := &mockCore{}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, countPercent(tt.args.s), "countPercent(%v)", tt.args.s)
		})
	}
}

func TestLogger_check(t *testing.T) {

}
