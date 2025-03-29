package glog

import (
	"go.uber.org/zap"
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

type LevelEnablerFunc = zap.LevelEnablerFunc

type AtomicLevel = zap.AtomicLevel

func ParseLevel(text string) (Level, error) {
	return zapcore.ParseLevel(text)
}
