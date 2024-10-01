//go:build go1.21

package slog

import (
	"context"
	"go.uber.org/zap"
	"log/slog"
	"testing"

	"github.com/ace-zhaoy/glog"
	"go.uber.org/zap/zapcore"
)

func TestHandler_Enabled(t *testing.T) {
	logger, _ := glog.NewDefault()
	handler := NewHandler(logger)

	tests := []struct {
		level    slog.Level
		expected bool
	}{
		{slog.LevelDebug, true},
		{slog.LevelInfo, true},
		{slog.LevelWarn, true},
		{slog.LevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := handler.Enabled(context.Background(), tt.level); got != tt.expected {
				t.Errorf("Enabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	logger, _ := glog.NewDefault()
	handler := NewHandler(logger)

	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: "test message",
	}

	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("Handle() error = %v", err)
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	logger, _ := glog.NewDefault()
	handler := NewHandler(logger)

	attrs := []slog.Attr{
		{Key: "key1", Value: slog.StringValue("value1")},
		{Key: "key2", Value: slog.IntValue(2)},
	}

	newHandler := handler.WithAttrs(attrs)
	if newHandler == handler {
		t.Errorf("WithAttrs() did not return a new handler")
	}
}

func TestHandler_WithGroup(t *testing.T) {
	logger, _ := glog.NewDefault()
	handler := NewHandler(logger)

	groupName := "testGroup"
	newHandler := handler.WithGroup(groupName)
	if newHandler == handler {
		t.Errorf("WithGroup() did not return a new handler")
	}
}

func TestLevelConverter(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected glog.Level
	}{
		{slog.LevelDebug, glog.LevelDebug},
		{slog.LevelInfo, glog.LevelInfo},
		{slog.LevelWarn, glog.LevelWarn},
		{slog.LevelError, glog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := LevelConverter(tt.level); got != tt.expected {
				t.Errorf("LevelConverter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAttr2Field(t *testing.T) {
	tests := []struct {
		attr     slog.Attr
		expected zapcore.Field
	}{
		{slog.Attr{Key: "key", Value: slog.StringValue("value")}, zap.String("key", "value")},
		{slog.Attr{Key: "key", Value: slog.IntValue(1)}, zap.Int64("key", 1)},
	}

	for _, tt := range tests {
		t.Run(tt.attr.Key, func(t *testing.T) {
			if got := attr2Field(tt.attr); got != tt.expected {
				t.Errorf("Attr2Field() = %v, want %v", got, tt.expected)
			}
		})
	}
}
