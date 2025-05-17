package slogflags

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoggerForTest(w io.Writer, opts ...Option) *slog.Logger {
	testOpts := append(opts, WithWriter(w), WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{Key: "time", Value: slog.StringValue("fake-time")}
		} else if a.Key == "file" {
			return slog.Attr{Key: "file", Value: slog.StringValue("file.go")}
		} else if a.Key == "function" {
			return slog.Attr{Key: "function", Value: slog.StringValue("github.com/csmith/slogflags.Test")}
		} else if a.Key == "line" {
			return slog.Attr{Key: "line", Value: slog.IntValue(87)}
		} else {
			return a
		}
	}))

	return Logger(testOpts...)
}

func Test_DefaultsToTextLogger(t *testing.T) {
	_ = flag.Set("log.format", "")

	w := new(bytes.Buffer)
	l := LoggerForTest(w)
	l.Warn("Test", "arg1", "arg2")

	assert.Equal(t, "time=fake-time level=WARN msg=Test arg1=arg2\n", w.String())
}

func Test_SetFormatToJson(t *testing.T) {
	_ = flag.Set("log.format", "json")

	w := new(bytes.Buffer)
	l := LoggerForTest(w)
	l.Warn("Test", "arg1", "arg2")

	assert.JSONEq(t, `{"arg1": "arg2", "level": "WARN", "msg": "Test", "time": "fake-time"}`, w.String())
}

func Test_AddingSource(t *testing.T) {
	_ = flag.Set("log.format", "json")

	w := new(bytes.Buffer)
	l := LoggerForTest(w, WithAddSource(true))
	l.Warn("Test", "arg1", "arg2")

	assert.JSONEq(t, `{
		"arg1": "arg2",
		"level": "WARN",
		"msg": "Test",
		"time": "fake-time",
		"source": {
			"file": "file.go",
			"function": "github.com/csmith/slogflags.Test",
			"line": 87
		}
	}`, w.String())
}

func Test_DefaultsToInfoLevel(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "")

	w := new(bytes.Buffer)
	l := LoggerForTest(w)
	l.Debug("Test")
	l.Info("Test")
	l.Warn("Test")
	l.Error("Test")

	assert.Equal(t, "time=fake-time level=INFO msg=Test\ntime=fake-time level=WARN msg=Test\ntime=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_WithCustomDefaultLevel(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "")

	w := new(bytes.Buffer)
	l := LoggerForTest(w, WithDefaultLogLevel(slog.LevelWarn))
	l.Warn("Test")
	l.Error("Test")

	assert.Equal(t, "time=fake-time level=WARN msg=Test\ntime=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_SetsBuiltInLevel(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "error")

	w := new(bytes.Buffer)
	l := LoggerForTest(w)
	l.Debug("Test")
	l.Info("Test")
	l.Warn("Test")
	l.Error("Test")

	assert.Equal(t, "time=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_SetsCustomLevel(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "shrug")

	custom := slog.Level(6)
	w := new(bytes.Buffer)
	l := LoggerForTest(w, WithCustomLevels(map[string]slog.Level{"shrug": custom}))
	l.Debug("Test")
	l.Info("Test")
	l.Warn("Test")
	l.Log(context.Background(), custom, "Test")
	l.Error("Test")

	assert.Equal(t, "time=fake-time level=SHRUG msg=Test\ntime=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_SetsDefault(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "error")

	w := new(bytes.Buffer)
	_ = LoggerForTest(w, WithSetDefault(true))

	slog.Error("Test")

	assert.Equal(t, "time=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_SetsOldLogLevel(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "error")

	w := new(bytes.Buffer)
	_ = LoggerForTest(w, WithSetDefault(true), WithOldLogLevel(slog.LevelError))

	log.Printf("Test")

	assert.Equal(t, "time=fake-time level=ERROR msg=Test\n", w.String())
}

func Test_CustomLevelWithReplaceAttr(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "shrug")

	custom := slog.Level(6)
	w := new(bytes.Buffer)
	
	// Test that both custom level handling and user ReplaceAttr work together
	replaceAttrCalled := false
	l := Logger(
		WithWriter(w),
		WithCustomLevels(map[string]slog.Level{"shrug": custom}),
		WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
			replaceAttrCalled = true
			if a.Key == slog.LevelKey {
				// The custom level should already be replaced with "SHRUG"
				assert.Equal(t, "SHRUG", a.Value.String())
			}
			if a.Key == "time" {
				return slog.Attr{Key: "time", Value: slog.StringValue("fake-time")}
			}
			return a
		}),
	)
	
	l.Log(context.Background(), custom, "Test")
	
	assert.True(t, replaceAttrCalled, "ReplaceAttr function should have been called")
	assert.Equal(t, "time=fake-time level=SHRUG msg=Test\n", w.String())
}

func Test_CustomLevelJSON(t *testing.T) {
	_ = flag.Set("log.format", "json")
	_ = flag.Set("log.level", "shrug")

	custom := slog.Level(6)
	w := new(bytes.Buffer)
	l := LoggerForTest(w, WithCustomLevels(map[string]slog.Level{"shrug": custom}))
	
	l.Log(context.Background(), custom, "Test", "arg1", "arg2")
	
	assert.JSONEq(t, `{
		"arg1": "arg2",
		"level": "SHRUG",
		"msg": "Test",
		"time": "fake-time"
	}`, w.String())
}

func Test_CustomLevelAliasLexicographicOrder(t *testing.T) {
	_ = flag.Set("log.format", "")
	_ = flag.Set("log.level", "custom")

	custom := slog.Level(10)
	w := new(bytes.Buffer)
	
	// First add "zebra", then "alpha" - should keep "alpha" since it's lexicographically first
	l := LoggerForTest(w, 
		WithCustomLevels(map[string]slog.Level{"zebra": custom}),
		WithCustomLevels(map[string]slog.Level{"alpha": custom}),
	)
	
	l.Log(context.Background(), custom, "Test")
	
	assert.Equal(t, "time=fake-time level=ALPHA msg=Test\n", w.String())
	
	// Now test the reverse order - add "alpha" then "zebra", should keep "alpha"
	w.Reset()
	l = LoggerForTest(w,
		WithCustomLevels(map[string]slog.Level{"alpha": custom}),
		WithCustomLevels(map[string]slog.Level{"zebra": custom}),
	)
	
	l.Log(context.Background(), custom, "Test")
	
	assert.Equal(t, "time=fake-time level=ALPHA msg=Test\n", w.String())
}
