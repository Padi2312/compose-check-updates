package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// CustomHandler is a slog.Handler that adds colorization based on log levels.
type CustomHandler struct {
	level  slog.Leveler
	output *os.File
}

// NewCustomHandler creates a new CustomHandler with the specified minimum log level.
func NewCustomHandler(level slog.Leveler, output *os.File) *CustomHandler {
	return &CustomHandler{
		level:  level,
		output: output,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle formats and writes the Record to the output.
func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	levelStr, color := getLevelStringAndColor(r.Level)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := r.Message
	reset := "\033[0m"

	// Collect attributes
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	// Combine message and attributes
	if len(attrs) > 0 {
		message = fmt.Sprintf("%s %s", message, attrsToString(attrs))
	}

	fmt.Fprintf(h.output, "%s[%s] %s %s%s\n", color, levelStr, timestamp, message, reset)
	return nil
}

// WithAttrs returns a new handler with the given attributes.
func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with the additional attributes
	newHandler := *h
	return &newHandler
}

// WithGroup returns a new handler with the given group.
func (h *CustomHandler) WithGroup(name string) slog.Handler {
	// Create a new handler with the group name
	newHandler := *h
	return &newHandler
}

// getLevelStringAndColor returns the string representation and color for a given log level.
func getLevelStringAndColor(level slog.Level) (string, string) {
	switch level {
	case slog.LevelDebug:
		return "DEBUG", "\033[34m" // Blue
	case slog.LevelInfo:
		return "INFO", "\033[32m" // Green
	case slog.LevelWarn:
		return "WARN", "\033[33m" // Yellow
	case slog.LevelError:
		return "ERROR", "\033[31m" // Red
	default:
		return "INFO", "\033[32m" // Default to INFO level
	}
}

// attrsToString formats the attributes slice into a string.
func attrsToString(attrs []string) string {
	return fmt.Sprintf("{%s}", join(attrs, ", "))
}

// join is a helper function to concatenate strings with a separator.
func join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
