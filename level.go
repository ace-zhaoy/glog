package glog

import (
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	LevelDebug = zapcore.DebugLevel
	LevelInfo  = zapcore.InfoLevel
	LevelWarn  = zapcore.WarnLevel
	LevelError = zapcore.ErrorLevel
)

type LevelEnabler = zapcore.LevelEnabler
