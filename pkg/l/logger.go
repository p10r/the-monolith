package l

import (
	slogtelegram "github.com/samber/slog-telegram/v2"
	"log/slog"
	"os"
)

func NewTextLogger() (log *slog.Logger) {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func NewAppLogger(h slog.Handler, app string) (log *slog.Logger) {
	return slog.New(h).With(slog.String("app", app))
}

func NewTelegramLogger(token, username, app string) (log *slog.Logger) {
	return slog.New(slogtelegram.Option{
		Level:    slog.LevelInfo,
		Token:    token,
		Username: username,
	}.NewTelegramHandler()).With(slog.String("app", app))
}

func Error(msg string, err error) (string, slog.Attr) {
	return msg, slog.Any("error", err)
}
