package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	ctxpkg "core/context"
)

// Logger is an enterprise-grade wrapper over slog.Logger with optimized context integration.
type Logger struct {
	handler slog.Handler
	attrs   []slog.Attr
	group   string
}

// Config holds logger configuration options.
type Config struct {
	Level      slog.Level
	AddSource  bool
	TimeFormat string // Optional: custom time format for text handlers
}

// DefaultConfig returns sensible defaults for production logging.
func DefaultConfig() *Config {
	return &Config{
		Level:     slog.LevelInfo,
		AddSource: false,
	}
}

// New creates a Logger from a slog.Handler.
func New(handler slog.Handler) *Logger {
	if handler == nil {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	return &Logger{handler: handler}
}

// NewJSON creates a JSON logger with the given configuration.
func NewJSON(w io.Writer, config *Config) *Logger {
	if w == nil {
		w = os.Stderr
	}
	if config == nil {
		config = DefaultConfig()
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	return New(slog.NewJSONHandler(w, opts))
}

// NewText creates a text logger with the given configuration.
func NewText(w io.Writer, config *Config) *Logger {
	if w == nil {
		w = os.Stderr
	}
	if config == nil {
		config = DefaultConfig()
	}

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	return New(slog.NewTextHandler(w, opts))
}

// With returns a new logger with the given attributes permanently attached.
func (l *Logger) With(args ...any) *Logger {
	if l == nil {
		return nil
	}

	attrs := argsToAttrs(args)
	return &Logger{
		handler: l.handler,
		attrs:   append(l.attrs, attrs...),
		group:   l.group,
	}
}

// WithGroup returns a new logger with the given group name.
func (l *Logger) WithGroup(name string) *Logger {
	if l == nil || name == "" {
		return l
	}

	newGroup := name
	if l.group != "" {
		newGroup = l.group + "." + name
	}

	return &Logger{
		handler: l.handler,
		attrs:   l.attrs,
		group:   newGroup,
	}
}

// LogAttrs logs a message with structured attributes at the given level.
func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if l == nil || !l.handler.Enabled(ctx, level) {
		return
	}

	record := slog.NewRecord(time.Now(), level, msg, 0)
	record.AddAttrs(l.attrs...)

	// Add context fields efficiently
	if ctx != nil {
		for k, v := range ctxpkg.Fields(ctx) {
			record.AddAttrs(slog.Any(k, v))
		}
	}

	// Add group prefix to attributes if needed
	if l.group != "" {
		for i := range attrs {
			attrs[i].Key = l.group + "." + attrs[i].Key
		}
	}

	record.AddAttrs(attrs...)
	_ = l.handler.Handle(ctx, record)
}

// Log logs a message with key-value pairs at the given level.
func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l == nil || !l.handler.Enabled(ctx, level) {
		return
	}

	attrs := argsToAttrs(args)
	l.LogAttrs(ctx, level, msg, attrs...)
}

// Debug logs at debug level.
func (l *Logger) Debug(msg string, args ...any) {
	l.Log(context.Background(), slog.LevelDebug, msg, args...)
}

// DebugContext logs at debug level with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at info level.
func (l *Logger) Info(msg string, args ...any) {
	l.Log(context.Background(), slog.LevelInfo, msg, args...)
}

// InfoContext logs at info level with context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string, args ...any) {
	l.Log(context.Background(), slog.LevelWarn, msg, args...)
}

// WarnContext logs at warn level with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at error level.
func (l *Logger) Error(msg string, args ...any) {
	l.Log(context.Background(), slog.LevelError, msg, args...)
}

// ErrorContext logs at error level with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, slog.LevelError, msg, args...)
}

// Enabled reports whether the logger handles records at the given level.
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	if l == nil {
		return false
	}
	return l.handler.Enabled(ctx, level)
}

// Handler returns the underlying slog.Handler.
func (l *Logger) Handler() slog.Handler {
	if l == nil {
		return nil
	}
	return l.handler
}

// ToSlog returns a *slog.Logger that uses this logger's handler.
// This allows integration with APIs that expect *slog.Logger.
func (l *Logger) ToSlog() *slog.Logger {
	if l == nil {
		return slog.Default()
	}

	handler := l.handler
	if len(l.attrs) > 0 {
		handler = handler.WithAttrs(l.attrs)
	}
	if l.group != "" {
		handler = handler.WithGroup(l.group)
	}

	return slog.New(handler)
}

// Global default logger management with proper synchronization
var (
	globalMu      sync.RWMutex
	defaultLogger *Logger
)

// SetDefault sets the global default logger.
// This should typically be called once during application startup.
func SetDefault(logger *Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	defaultLogger = logger
}

// Default returns the global default logger.
// If no logger has been set, it returns a JSON logger writing to stderr.
func Default() *Logger {
	globalMu.RLock()
	if defaultLogger != nil {
		logger := defaultLogger
		globalMu.RUnlock()
		return logger
	}
	globalMu.RUnlock()

	// Initialize default logger if not set
	globalMu.Lock()
	defer globalMu.Unlock()

	// Double-check after acquiring write lock
	if defaultLogger == nil {
		defaultLogger = NewJSON(os.Stderr, DefaultConfig())
	}

	return defaultLogger
}

// Package-level convenience functions that use the default logger

// Debug logs at debug level using the default logger.
func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

// DebugContext logs at debug level with context using the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Default().DebugContext(ctx, msg, args...)
}

// Info logs at info level using the default logger.
func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

// InfoContext logs at info level with context using the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Default().InfoContext(ctx, msg, args...)
}

// Warn logs at warn level using the default logger.
func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

// WarnContext logs at warn level with context using the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Default().WarnContext(ctx, msg, args...)
}

// Error logs at error level using the default logger.
func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

// ErrorContext logs at error level with context using the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default().ErrorContext(ctx, msg, args...)
}

// Helper functions

func argsToAttrs(args []any) []slog.Attr {
	if len(args) == 0 {
		return nil
	}

	attrs := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		attrs = append(attrs, slog.Any(key, args[i+1]))
	}

	return attrs
}
