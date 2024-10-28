package log

import (
	"log"
	"log/slog"
	"os"
)

func InitSlog(level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	h := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(h)
}

func InitLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
