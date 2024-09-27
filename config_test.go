package glog

import (
	"testing"
)

func TestConfig_Build(t *testing.T) {
	cfg := NewDefaultConfig()
	logger, err := cfg.Build()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if logger == nil {
		t.Error("Expected logger to be created, but got nil")
	}

	stackLevel := LevelError
	cfg = &Config{
		Level:      LevelInfo,
		AddCaller:  true,
		StackLevel: &stackLevel,
		Sampling: &SamplingConfig{
			Initial:    50,
			Thereafter: 50,
		},
		Core: CoreConfig{
			Encoding: "json",
		},
	}
	logger, err = cfg.Build()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if logger == nil {
		t.Error("Expected logger to be created, but got nil")
	}
}

func TestConfig_buildOptions(t *testing.T) {
	cfg := &Config{}
	opts := cfg.buildOptions()
	if len(opts) != 0 {
		t.Error("Expected no options, but got", opts)
	}

	cfg = &Config{
		InitialFields: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}
	opts = cfg.buildOptions()
	if len(opts) != 1 {
		t.Error("Expected one option, but got", opts)
	}
}
