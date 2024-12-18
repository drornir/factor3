package log

import (
	"context"
	"sync"
)

var (
	globalLogger     Logger
	globalLoggerLock sync.RWMutex
)

func init() {
	globalLogger = NoopLogger{}
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
