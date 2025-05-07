/*
Package slogflags provides flags to configure [log/slog].

# Basic usage

Simply call [flag.Parse] and then call [Logger] to obtain a configured slog
instance. Two new flags will be available to users of your app: `--log.level`
which accepts a textual level ("debug", "info", "warn" or "error") and
`--log.format` which accepts either "text" or "json".

	flag.Parse()
	logger := slogflags.Logger()
	logger.Warn("This is not a drill", "key", "value", "etc", "etc)

# Custom levels

If you define your own log levels, you can pass them to [Logger] using
[WithCustomLevels]. Users can then specify them in the `log.level` flag.

# Setting as the default logger

Pass [WithSetDefault](true) when calling [Logger] to register the new instance
as the default logger for [log/slog]. You can then call [log/slog.Warn] etc
directly.

# Redirecting old log calls

If your code or libraries use [log] rather than [log/slog] you can redirect them
by setting a default logger as above. All [log] calls will be given the same
log level, which you can alter using [WithOldLogLevel]. e.g.:

	flag.Parse()
	_ = slogflags.Logger(slogflags.WithSetDefault(true), slogflags.WithOldLogLevel(slog.LevelWarn))
	log.Printf("hi")
	// Prints: time=... level=WARN msg=hi

# Other advanced usage

You can customise other behaviour of the created logger using
[WithDefaultLogLevel], [WithWriter], [WithAddSource] and [WithReplaceAttr].
See the documentation for those funcs for more details.
*/
package slogflags
