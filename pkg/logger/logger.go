package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	White  = "\033[97m"
)

// ColorHandler wraps a slog.Handler to add colored output
type ColorHandler struct {
	handler slog.Handler
	writer  io.Writer
	mu      sync.Mutex
}

func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	return &ColorHandler{
		handler: slog.NewTextHandler(w, opts),
		writer:  w,
	}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format level with color
	levelColor := White
	switch r.Level {
	case slog.LevelDebug:
		levelColor = Gray
	case slog.LevelInfo:
		levelColor = White
	case slog.LevelWarn:
		levelColor = Yellow
	case slog.LevelError:
		levelColor = Red
	}

	// Apply color to entire line
	fmt.Fprint(h.writer, levelColor)

	// Extract tag and other attributes first to calculate full prefix length
	var tag string
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "tag" {
			tag = a.Value.String()
		} else {
			attrs = append(attrs, a)
		}
		return true
	})

	// Build the line: [time level]: [tag] message
	timeStr := r.Time.Format("15:04:05.000")

	var line string
	if tag != "" {
		line = fmt.Sprintf("[%s %s]: [%s] %s", timeStr, r.Level.String(), tag, r.Message)
	} else {
		line = fmt.Sprintf("[%s %s]: %s", timeStr, r.Level.String(), r.Message)
	}

	// Write the line with padding if there are attributes
	if len(attrs) > 0 {
		// Format line to 100 columns
		fmt.Fprintf(h.writer, "%-100s", line)

		// Reset color before writing attributes
		fmt.Fprint(h.writer, Reset)

		// Write attributes with green keys
		for i, a := range attrs {
			if i > 0 {
				fmt.Fprint(h.writer, " ")
			}
			// Apply green color to keys
			fmt.Fprint(h.writer, Green)
			fmt.Fprintf(h.writer, "%s", a.Key)
			fmt.Fprint(h.writer, Gray)

			// Only add quotes for string values and errors
			switch a.Value.Kind() {
			case slog.KindString:
				fmt.Fprintf(h.writer, "=\"%s\"", a.Value.String())
			default:
				if _, ok := a.Value.Any().(error); ok {
					fmt.Fprintf(h.writer, "=\"%v\"", a.Value)
				} else {
					fmt.Fprintf(h.writer, "=%v", a.Value)
				}
			}
		}
	} else {
		// No attributes, just write the line
		fmt.Fprint(h.writer, line)
	}

	// Reset color at end of line
	fmt.Fprintf(h.writer, "%s\n", Reset)
	return nil
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	return &ColorHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := NewColorHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
}

// SetLevel sets the minimum log level
func SetLevel(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := NewColorHandler(os.Stdout, opts)
	defaultLogger = slog.New(handler)
}

// Info logs an info message with a tag and optional key-value pairs
// Usage: logger.Info("TAG", "message", "key1", value1, "key2", value2)
func Info(tag string, msg string, args ...interface{}) {
	attrs := make([]interface{}, 0, len(args)+2)
	attrs = append(attrs, "tag", tag)
	attrs = append(attrs, args...)
	defaultLogger.Info(msg, attrs...)
}

// Debug logs a debug message with a tag and optional key-value pairs
func Debug(tag string, msg string, args ...interface{}) {
	attrs := make([]interface{}, 0, len(args)+2)
	attrs = append(attrs, "tag", tag)
	attrs = append(attrs, args...)
	defaultLogger.Debug(msg, attrs...)
}

// Warn logs a warning message with a tag and optional key-value pairs
func Warn(tag string, msg string, args ...interface{}) {
	attrs := make([]interface{}, 0, len(args)+2)
	attrs = append(attrs, "tag", tag)
	attrs = append(attrs, args...)
	defaultLogger.Warn(msg, attrs...)
}

// Error logs an error message with a tag and optional key-value pairs
func Error(tag string, msg string, args ...interface{}) {
	attrs := make([]interface{}, 0, len(args)+2)
	attrs = append(attrs, "tag", tag)
	attrs = append(attrs, args...)
	defaultLogger.Error(msg, attrs...)
}

// Fatal logs an error message with a tag and exits
func Fatal(tag string, msg string, args ...interface{}) {
	attrs := make([]interface{}, 0, len(args)+2)
	attrs = append(attrs, "tag", tag)
	attrs = append(attrs, args...)
	defaultLogger.Error(msg, attrs...)
	os.Exit(1)
}

// Legacy compatibility functions

func Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	defaultLogger.Info(msg)
}

func Println(msg string, args ...interface{}) {
	if len(args) > 0 {
		defaultLogger.Info(msg, args...)
	} else {
		defaultLogger.Info(msg)
	}
}

func Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	defaultLogger.Debug(msg)
}

func Debugln(msg string) {
	defaultLogger.Debug(msg)
}

func Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	defaultLogger.Warn(msg)
}

func Warnln(msg string, args ...interface{}) {
	defaultLogger.Warn(msg)
}

func Errorln(v ...interface{}) {
	if len(v) == 0 {
		defaultLogger.Error("")
		return
	}

	// If only one argument, treat as simple message
	if len(v) == 1 {
		defaultLogger.Error(fmt.Sprint(v[0]))
		return
	}

	// Multiple arguments - first is message, rest are key-value pairs
	msg := fmt.Sprint(v[0])
	args := v[1:]

	// If odd number of remaining args, just append them to message
	if len(args)%2 != 0 {
		msg = fmt.Sprintln(v...)
		if len(msg) > 0 && msg[len(msg)-1] == '\n' {
			msg = msg[:len(msg)-1]
		}
		defaultLogger.Error(msg)
		return
	}

	// Even number of remaining args - treat as key-value pairs
	defaultLogger.Error(msg, args...)
}

func Fatalln(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	// Remove trailing newline that Sprintln adds
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	defaultLogger.Error(msg)
	os.Exit(1)
}
