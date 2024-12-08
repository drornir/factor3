package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func SlogReplacerMinimal() func(groups []string, a slog.Attr) slog.Attr {
	initialTime := time.Now()
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			if a.Value.Kind() != slog.KindTime {
				return a
			}
			orig := a.Value.Time()
			dur := orig.Sub(initialTime)
			a.Value = slog.DurationValue(dur.Round(time.Microsecond))
			return a
		}
		if a.Key == slog.SourceKey {
			src, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}
			cwd, err := os.Getwd()
			cobra.CheckErr(err)
			gopath := os.Getenv("GOPATH")
			fmt.Println("src", src)
			for _, prefix := range []string{gopath, cwd} {
				if smaller, ok := strings.CutPrefix(src.File, prefix); ok {
					src.File = smaller
				}
			}
		}
		return a
	}
}

type PflagLeveler struct{ Flag *string }

var _ slog.Leveler = PflagLeveler{}

func (l PflagLeveler) Level() slog.Level {
	if l.Flag == nil {
		return slog.LevelInfo
	}
	s := *l.Flag
	return ParseLevel(s)
}
