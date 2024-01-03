package log

import (
	"log/slog"
	"os"
)

var Logger = newLogger()

func newLogger() *slog.Logger {
	opt := slog.HandlerOptions{}
	return slog.New(slog.NewJSONHandler(os.Stdout, &opt))
}
