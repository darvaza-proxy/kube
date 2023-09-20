package main

import (
	"fmt"

	"darvaza.org/sidecar/pkg/logger/zap"
	"darvaza.org/slog"
)

var (
	verbosity int
)

// fatal is a convenience wrapper for slog.Logger.Fatal().Print()
func fatal(err error, msg string, args ...any) {
	l := newLogger().Fatal()
	if err != nil {
		l = l.WithField(slog.ErrorFieldName, err)
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.Print(msg)

	panic("unreachable")
}

func newLogger() slog.Logger {
	level := slog.Error + slog.LogLevel(verbosity)

	switch {
	case level < slog.Fatal:
		level = slog.Fatal
	case level > slog.Debug:
		level = slog.Debug
	}

	return zap.New(level)
}

func init() {
	pflags := svc.PersistentFlags()
	pflags.CountVarP(&verbosity, "verbosity", "v", "increase the verbosity level to Warn, Info or Debug")
}
