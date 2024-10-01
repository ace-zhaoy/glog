package cores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type mockCore struct {
	enabled bool
	fields  []zapcore.Field
}

func (m *mockCore) Enabled(lvl zapcore.Level) bool {
	return m.enabled
}

func (m *mockCore) With(fields []zapcore.Field) zapcore.Core {
	return &mockCore{
		enabled: m.enabled,
		fields:  append(m.fields, fields...),
	}
}

func (m *mockCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce
}

func (m *mockCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	m.fields = append(m.fields, fields...)
	return nil
}

func (m *mockCore) Sync() error {
	return nil
}

func TestNewLazyCore(t *testing.T) {
	core := &mockCore{}
	lazyCore := NewLazyCore(core)
	assert.NotNil(t, lazyCore, "Expected NewLazyCore to return a non-nil LazyCore")
}

func TestLazyCore_Enabled(t *testing.T) {
	core := &mockCore{enabled: true}
	lazyCore := NewLazyCore(core)
	assert.True(t, lazyCore.Enabled(zapcore.DebugLevel), "Expected Enabled to return true")
}

func TestLazyCore_With(t *testing.T) {
	core := &mockCore{}
	lazyCore := NewLazyCore(core)
	newCore := lazyCore.With([]zapcore.Field{zapcore.Field{Key: "key", Type: zapcore.StringType, String: "value"}})
	assert.NotNil(t, newCore, "Expected With to return a non-nil Core")
}

func TestLazyCore_Check(t *testing.T) {
	core := &mockCore{}
	lazyCore := NewLazyCore(core)
	entry := zapcore.Entry{}
	checkedEntry := &zapcore.CheckedEntry{}
	result := lazyCore.Check(entry, checkedEntry)
	assert.Equal(t, checkedEntry, result, "Expected Check to return the same CheckedEntry")
}

func TestLazyCore_Write(t *testing.T) {
	core := &mockCore{}
	lazyCore := NewLazyCore(core)
	entry := zapcore.Entry{}
	fields := []zapcore.Field{zapcore.Field{Key: "key", Type: zapcore.StringType, String: "value"}}
	err := lazyCore.Write(entry, fields)
	assert.NoError(t, err, "Expected Write to not return an error")
	assert.Equal(t, fields, core.fields, "Expected fields to be written to the core")
}

func TestLazyCore_Sync(t *testing.T) {
	core := &mockCore{}
	lazyCore := NewLazyCore(core)
	err := lazyCore.Sync()
	assert.NoError(t, err, "Expected Sync to not return an error")
}
