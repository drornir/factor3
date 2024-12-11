package tests

import (
	"log/slog"
	"os"
	"sync"

	"github.com/drornir/factor3/pkg/log"
)

var initOnce = sync.OnceFunc(func() {
	initLogger()
})

func init() { Init() }

func Init() {
	initOnce()
}

func initLogger() {
	logOut := os.Stdout // TODO  configurable

	var sloggerHandler slog.Handler
	sloggerOpts := &slog.HandlerOptions{
		// AddSource: log.ParseLevel(flagLogLevel) <= slog.LevelDebug,
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: log.SlogReplacerMinimal(),
	}

	sloggerHandler = slog.NewTextHandler(logOut, sloggerOpts)

	slogger := slog.New(sloggerHandler)
	log.SetGlobal(log.WrapSlogger(slogger))
}
