package glog

import "github.com/ace-zhaoy/glog/cores"

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WrapCore(f func(core cores.Core) cores.Core) Option {
	return optionFunc(func(log *Logger) {
		log.core = f(log.core)
	})
}

func WithName(name string) Option {
	return optionFunc(func(l *Logger) {
		l.name = name
	})
}

func WithCaller(enabled bool) Option {
	return optionFunc(func(log *Logger) {
		log.addCaller = enabled
	})
}

func AddCaller() Option {
	return WithCaller(true)
}

func WithStack(lvl LevelEnabler) Option {
	return optionFunc(func(l *Logger) {
		l.addStack = lvl
	})
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *Logger) {
		l.callerSkip = skip
	})
}

func WithFormat(formatEnabled bool) Option {
	return optionFunc(func(l *Logger) {
		l.formatEnabled = formatEnabled
	})
}

func WithContextHandlers(handlers ...ContextHandler) Option {
	return optionFunc(func(l *Logger) {
		l.contextHandlers = handlers
	})
}

func AddContextHandlers(handlers ...ContextHandler) Option {
	return optionFunc(func(l *Logger) {
		l.contextHandlers = append(l.contextHandlers, handlers...)
	})
}
