package log

import (
	"context"
	"log/slog"
)

const SlogLevelTrace slog.Level = slog.LevelDebug - 4

func WrapSlogger(slogger *slog.Logger) Logger {
	return &sloggerWrapper{
		sl: slogger,
	}
}

type sloggerWrapper struct {
	sl *slog.Logger
}

func (l *sloggerWrapper) T(ctx context.Context, msg string, a ...any) {
	l.log(ctx, SlogLevelTrace, msg, a...)
}

func (l *sloggerWrapper) D(ctx context.Context, msg string, a ...any) {
	l.log(ctx, slog.LevelDebug, msg, a...)
}

func (l *sloggerWrapper) I(ctx context.Context, msg string, a ...any) {
	l.log(ctx, slog.LevelInfo, msg, a...)
}

func (l *sloggerWrapper) W(ctx context.Context, msg string, a ...any) {
	l.log(ctx, slog.LevelWarn, msg, a...)
}

func (l *sloggerWrapper) E(ctx context.Context, msg string, a ...any) {
	l.log(ctx, slog.LevelError, msg, a...)
}

func (l *sloggerWrapper) Slogger() *slog.Logger {
	return l.sl
}

func (l *sloggerWrapper) log(ctx context.Context, level slog.Level, msg string, a ...any) {
	l.sl.Log(ctx, level, msg, a...)
}
