package glog

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestCoreConfig_buildEncoder(t *testing.T) {
	cfg := CoreConfig{
		Encoding:      "json",
		EncoderConfig: zapcore.EncoderConfig{},
	}

	enc, err := cfg.buildEncoder()
	assert.NoError(t, err, "Expected no error when building JSON encoder")
	assert.IsType(t, zapcore.NewJSONEncoder(zapcore.EncoderConfig{}), enc, "Expected JSON encoder")

	cfg.Encoding = "console"
	enc, err = cfg.buildEncoder()
	assert.NoError(t, err, "Expected no error when building Console encoder")
	assert.IsType(t, zapcore.NewConsoleEncoder(zapcore.EncoderConfig{}), enc, "Expected Console encoder")

	cfg.Encoding = "unsupported"
	enc, err = cfg.buildEncoder()
	assert.Error(t, err, "Expected error for unsupported encoding")
}

func TestCoreConfig_openSinks(t *testing.T) {
	cfg := CoreConfig{
		OutputPaths: []string{"stdout"},
	}

	sink, err := cfg.openSinks()
	assert.NoError(t, err, "Expected no error when opening stdout sink")
	assert.NotNil(t, sink, "Expected valid WriteSyncer sink")

	cfg.OutputPaths = []string{""}
	sink, err = cfg.openSinks()
	assert.Error(t, err, "Expected error when opening invalid path sink")
}

func TestCoreConfig_Build(t *testing.T) {
	lvl := zapcore.InfoLevel
	cfg := CoreConfig{
		Encoding:      "json",
		EncoderConfig: zapcore.EncoderConfig{},
		OutputPaths:   []string{"stdout"},
	}

	core, err := cfg.Build(lvl)
	assert.NoError(t, err, "Expected no error when building Core")
	assert.NotNil(t, core, "Expected valid Core")

	cfg.Encoding = "unsupported"
	core, err = cfg.Build(lvl)
	assert.Error(t, err, "Expected error when building Core with unsupported encoding")
}
