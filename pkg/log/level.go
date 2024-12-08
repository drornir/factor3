package log

import (
	"log/slog"
	"strings"
)

// ParseLevel is a more forgiving version for parsing a string into an slog.Level
func ParseLevel(s string) slog.Level {
	s = strings.TrimSpace(s)
	if s == "" {
		s = "info"
	}
	var logLevel slog.Level
	switch strings.ToLower(s)[0] {
	case 't':
		logLevel = SlogLevelTrace
	case 'd':
		logLevel = slog.LevelDebug
	case 'i':
		logLevel = slog.LevelInfo
	case 'w':
		logLevel = slog.LevelWarn
	case 'e':
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return logLevel
}
