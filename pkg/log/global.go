package log

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	globalLogger     Logger
	globalLoggerLock sync.RWMutex
)

func init() {
	// users should override this. But if they don't, let's put a sane default
	globalLogger = WrapSlogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
}

func GG() Logger {
	return GetGlobal()
}

func GetGlobal() Logger {
	globalLoggerLock.RLock()
	defer globalLoggerLock.RUnlock()
	return globalLogger
}

func SetGlobal(l Logger) {
	globalLoggerLock.Lock()
	defer globalLoggerLock.Unlock()
	globalLogger = l

	l.T(context.Background(), "log level set to 'trace'")
}
