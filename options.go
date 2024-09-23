package glog

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WrapCore(f func(core Core) Core) Option {
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
		l.stackLevel = lvl
	})
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *Logger) {
		l.callerSkip += skip
	})
}

func WithCallerSkip(skip int) Option {
	return optionFunc(func(l *Logger) {
		l.callerSkip = skip
	})
}

func WithFormat(formatEnabled bool) Option {
	return optionFunc(func(l *Logger) {
		l.formatEnabled = formatEnabled
	})
}

func WithFormatEnabled() Option {
	return WithFormat(true)
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

func WithFields(fields ...Field) Option {
	return optionFunc(func(l *Logger) {
		l.core = l.core.With(fields)
	})
}
