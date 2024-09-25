package glog

import (
	"github.com/ace-zhaoy/glog/cores"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sort"
	"time"
)

type SamplingConfig = zap.SamplingConfig

type Config struct {
	Name          string            `json:"name" yaml:"name"`
	Level         Level             `json:"level" yaml:"level"`
	LazyDisabled  bool              `json:"lazyDisabled" yaml:"lazyDisabled"`
	AddCaller     bool              `json:"addCaller" yaml:"addCaller"`
	StackLevel    *Level            `json:"stackLevel" yaml:"stackLevel"`
	CallerSkip    int               `json:"callerSkip" yaml:"callerSkip"`
	FormatEnabled bool              `json:"formatEnabled" yaml:"formatEnabled"`
	ContextFields map[string]string `json:"contextFields" yaml:"contextFields"`
	Sampling      *SamplingConfig   `json:"sampling" yaml:"sampling"`
	InitialFields map[string]any    `json:"initialFields" yaml:"initialFields"`
	Core          CoreConfig        `json:"core" yaml:"core"`
}

func (c *Config) buildOptions() []Option {
	opts := make([]Option, 0, 10)

	if c.Name != "" {
		opts = append(opts, WithName(c.Name))
	}

	if c.LazyDisabled {
		opts = append(opts, WrapCore(func(core Core) Core {
			return cores.NewLazyCore(core)
		}))
	}

	if c.AddCaller {
		opts = append(opts, AddCaller())
	}
	if c.StackLevel != nil {
		opts = append(opts, WithStack(c.StackLevel))
	}
	if c.CallerSkip != 0 {
		opts = append(opts, WithCallerSkip(c.CallerSkip))
	}

	if c.FormatEnabled {
		opts = append(opts, WithFormatEnabled())
	}

	if len(c.ContextFields) > 0 {
		keys := make([]string, 0, len(c.ContextFields))
		for k := range c.ContextFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		contextHandlers := make([]ContextHandler, 0, len(c.ContextFields))
		for _, k := range keys {
			contextHandlers = append(contextHandlers, BuildContextHandler(k, c.ContextFields[k]))
		}
		opts = append(opts, WithContextHandlers(contextHandlers...))
	}

	if len(c.InitialFields) > 0 {
		keys := make([]string, 0, len(c.InitialFields))
		for k := range c.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fds := make([]Field, 0, len(c.InitialFields))
		for _, k := range keys {
			fds = append(fds, Any(k, c.InitialFields[k]))
		}

		opts = append(opts, WrapCore(func(core Core) Core {
			return core.With(fds)
		}))
	}

	if c.Sampling != nil {
		opts = append(opts, WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if c.Sampling.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(c.Sampling.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				c.Sampling.Initial,
				c.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	return opts
}
func (c *Config) Build(opts ...Option) (*Logger, error) {
	core, err := c.Core.Build(c.Level)
	if err != nil {
		return nil, err
	}

	return NewLogger(core, c.buildOptions()...).WithOptions(opts...), nil
}

func NewDefaultConfig() *Config {
	stackLevel := LevelError
	return &Config{
		Level:      LevelDebug,
		AddCaller:  true,
		StackLevel: &stackLevel,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Core: CoreConfig{
			Encoding: "json",
			EncoderConfig: EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths: []string{"stderr"},
		},
	}
}
