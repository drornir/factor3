package log

import (
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
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

			var paths []string
			cwd, err := os.Getwd()
			if err == nil {
				paths = append(paths, cwd+"/")
			}
			gopath := os.Getenv("GOPATH")
			if gopath != "" {
				paths = append(paths, gopath+"/")
			}
			if modfile := os.Getenv("GOMOD"); modfile != "" {
				moddir := path.Dir(modfile)
				paths = append(paths, moddir+"/")
			}
			homedir, err := os.UserHomeDir()
			if err == nil {
				paths = append(paths, homedir+"/")
			}

			// fmt.Println("homedir=", homedir, "cwd=", cwd, "gopath=", gopath)
			for _, prefix := range paths {
				if smaller, ok := strings.CutPrefix(src.File, prefix); ok {
					src.File = smaller
				}
			}
			return a
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
