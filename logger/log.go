package logger

import (
	"log"
	"log/slog"
	"os"
)

func InitSlog(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	h := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(h))
}

func InitLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
