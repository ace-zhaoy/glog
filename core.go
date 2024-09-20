package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Core = zapcore.Core

type EncoderConfig = zapcore.EncoderConfig

type CoreConfig struct {
	Encoding      string        `json:"encoding" yaml:"encoding"`
	EncoderConfig EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	OutputPaths   []string      `json:"outputPaths" yaml:"outputPaths"`
}

func (c *CoreConfig) buildEncoder() (zapcore.Encoder, error) {
	switch c.Encoding {
	case "json":
		return zapcore.NewJSONEncoder(c.EncoderConfig), nil
	case "console":
		return zapcore.NewConsoleEncoder(c.EncoderConfig), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", c.Encoding)
	}
}

func (c *CoreConfig) openSinks() (zapcore.WriteSyncer, error) {
	sink, _, err := zap.Open(c.OutputPaths...)
	return sink, err
}

func (c *CoreConfig) Build(lvl LevelEnabler) (core Core, err error) {
	enc, err := c.buildEncoder()
	if err != nil {
		return
	}
	sink, err := c.openSinks()
	if err != nil {
		return
	}
	core = zapcore.NewCore(enc, sink, lvl)
	return
}
