package glog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

type mockCore struct{}

func NewMockCore() Core                                                                { return mockCore{} }
func (mockCore) Enabled(Level) bool                                                    { return false }
func (m mockCore) With(_ []Field) Core                                                 { return m }
func (mockCore) Check(_ zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry { return ce }
func (mockCore) Write(zapcore.Entry, []Field) error                                    { return nil }
func (mockCore) Sync() error                                                           { return nil }

func TestWrapCore(t *testing.T) {
	logger := &Logger{}
	option := WrapCore(func(core Core) Core {
		return mockCore{}
	})

	option.apply(logger)
	assert.NotNil(t, logger.core, "Expected core to be wrapped")
}

func TestWithName(t *testing.T) {
	logger := &Logger{}
	option := WithName("test_logger")

	option.apply(logger)
	assert.Equal(t, "test_logger", logger.name, "Expected logger name to be 'test_logger'")
}

func TestWithCaller(t *testing.T) {
	logger := &Logger{}
	option := WithCaller(true)

	option.apply(logger)
	assert.True(t, logger.addCaller, "Expected addCaller to be true")

	option = WithCaller(false)
	option.apply(logger)
	assert.False(t, logger.addCaller, "Expected addCaller to be false")
}

func TestAddCaller(t *testing.T) {
	logger := &Logger{}
	option := AddCaller()

	option.apply(logger)
	assert.True(t, logger.addCaller, "Expected addCaller to be true from AddCaller")
}

func TestWithStack(t *testing.T) {
	logger := &Logger{}
	lvl := LevelInfo
	option := WithStack(lvl)

	option.apply(logger)
	assert.Equal(t, lvl, logger.stackLevel, "Expected stackLevel to be set to InfoLevel")
}

func TestAddCallerSkip(t *testing.T) {
	logger := &Logger{}
	option := AddCallerSkip(2)

	option.apply(logger)
	assert.Equal(t, 2, logger.callerSkip, "Expected callerSkip to be incremented by 2")
}

func TestWithCallerSkip(t *testing.T) {
	logger := &Logger{}
	option := WithCallerSkip(3)

	option.apply(logger)
	assert.Equal(t, 3, logger.callerSkip, "Expected callerSkip to be set to 3")
}

func TestWithFormatEnabled(t *testing.T) {
	logger := &Logger{}
	option := WithFormatEnabled()

	option.apply(logger)
	assert.True(t, logger.formatEnabled, "Expected formatEnabled to be true")
}

func TestWithContextHandlers(t *testing.T) {
	logger := &Logger{}
	handler := func(ctx context.Context, record *Record) {}
	option := WithContextHandlers(handler)

	option.apply(logger)
	assert.NotNil(t, logger.contextHandlers, "Expected contextHandlers to be set")
	assert.Equal(t, 1, len(logger.contextHandlers), "Expected one contextHandler")
}

func TestAddContextHandlers(t *testing.T) {
	logger := &Logger{}
	handler1 := func(ctx context.Context, record *Record) {}
	handler2 := func(ctx context.Context, record *Record) {}
	option := AddContextHandlers(handler1)

	option.apply(logger)
	assert.Equal(t, 1, len(logger.contextHandlers), "Expected one contextHandler")

	option = AddContextHandlers(handler2)
	option.apply(logger)
	assert.Equal(t, 2, len(logger.contextHandlers), "Expected two contextHandlers")
}

func TestWithFields(t *testing.T) {
	logger := &Logger{core: mockCore{}}
	field := zapcore.Field{Key: "key", String: "value"}
	option := WithFields(field)

	option.apply(logger)
	assert.NotNil(t, logger.core, "Expected core to be set with fields")
}
