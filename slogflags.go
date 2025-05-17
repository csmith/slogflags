package slogflags

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"strings"
)

var (
	logLevel  = flag.String("log.level", "", "Lowest level of logs that should be output")
	logFormat = flag.String("log.format", "text", "Format of log output ('json' or 'text')")

	defaultLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

// Logger creates a new [log/slog.Logger] configured according to the options
// and flags.
//
// [flag.Parse] must be called prior to calling this method.
func Logger(opts ...Option) *slog.Logger {
	c := newConfig(opts)

	slog.SetLogLoggerLevel(c.oldLogLevel)

	var handlerOpts = &slog.HandlerOptions{
		AddSource:   c.addSource,
		Level:       c.level(*logLevel),
		ReplaceAttr: c.levelReplaceAttr,
	}

	var handler slog.Handler
	if *logFormat == "json" {
		handler = slog.NewJSONHandler(c.writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(c.writer, handlerOpts)
	}

	logger := slog.New(handler)
	if c.setDefault {
		slog.SetDefault(logger)
	}
	return logger
}

type config struct {
	addSource        bool
	customLevels     map[string]slog.Level
	customLevelNames map[slog.Level]string
	defaultLevel     slog.Level
	oldLogLevel      slog.Level
	replaceAttr      func(groups []string, a slog.Attr) slog.Attr
	setDefault       bool
	writer           io.Writer
}

func newConfig(opts []Option) *config {
	c := &config{
		addSource:        false,
		defaultLevel:     slog.LevelInfo,
		oldLogLevel:      slog.LevelInfo,
		customLevels:     map[string]slog.Level{},
		customLevelNames: map[slog.Level]string{},
		replaceAttr:      nil,
		setDefault:       false,
		writer:           os.Stdout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *config) level(requested string) slog.Level {
	target := strings.ToLower(requested)

	if r, ok := defaultLevels[target]; ok {
		return r
	}

	if r, ok := c.customLevels[target]; ok {
		return r
	}

	return c.defaultLevel
}

func (c *config) levelReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		if name, ok := c.customLevelNames[a.Value.Any().(slog.Level)]; ok {
			a = slog.String(slog.LevelKey, name)
		}
	}

	if c.replaceAttr != nil {
		return c.replaceAttr(groups, a)
	}

	return a
}

type Option func(*config)

// WithAddSource controls whether the source code location will be added to
// log lines. See [log/slog.HandlerOptions.AddSource].
func WithAddSource(addSource bool) Option {
	return func(c *config) {
		c.addSource = addSource
	}
}

// WithCustomLevels adds extra levels to the defaults available in the
// `log.level` flag. The same level may be specified with multiple different
// keys to provide aliases.
//
// The "level" attribute of log messages will automatically be rewritten
// if the level matches a custom attribute. If multiple aliases are defined
// for the same level, the alias that comes first lexicographically will
// be used.
func WithCustomLevels(levels map[string]slog.Level) Option {
	return func(c *config) {
		for k, v := range levels {
			c.customLevels[strings.ToLower(k)] = v

			upperK := strings.ToUpper(k)
			if existing, ok := c.customLevelNames[v]; !ok || upperK < existing {
				c.customLevelNames[v] = upperK
			}
		}
	}
}

// WithDefaultLogLevel sets the default level that will be used if the
// `log.level` flag is not set. If not provided, the default
// is [log/slog.LevelInfo].
func WithDefaultLogLevel(level slog.Level) Option {
	return func(c *config) {
		c.defaultLevel = level
	}
}

// WithOldLogLevel sets the level that should be used when interoping with the
// older [log] package. See [log/slog.SetLogLoggerLevel]. If not provided, the
// default is [log/slog.LevelInfo].
func WithOldLogLevel(level slog.Level) Option {
	return func(c *config) {
		c.oldLogLevel = level
	}
}

// WithReplaceAttr allows setting an attribute replacement func on the logger.
// This can be used to rewrite attribute names or values.
// See [log/slog.HandlerOptions.ReplaceAttr].
func WithReplaceAttr(fn func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(c *config) {
		c.replaceAttr = fn
	}
}

// WithSetDefault sets whether the logger should be set as the default [log/slog]
// logger. See [log/slog.SetDefault].
func WithSetDefault(setDefault bool) Option {
	return func(c *config) {
		c.setDefault = setDefault
	}
}

// WithWriter sets a custom writer to be used for the log output. Defaults to
// [os.Stdout].
func WithWriter(w io.Writer) Option {
	return func(c *config) {
		c.writer = w
	}
}
