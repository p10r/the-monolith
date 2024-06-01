package logging

import (
	"log/slog"
	"os"
)

func NewTextLogger() (log *slog.Logger) {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func NewAppLogger(h slog.Handler, app string) (log *slog.Logger) {
	return slog.New(h).With(slog.String("app", app))
}
