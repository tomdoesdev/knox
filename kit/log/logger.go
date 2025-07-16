package log

import (
	"log/slog"
	"os"
)

type Level string
type LoggerType string

const (
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

const (
	Text LoggerType = "text"
	JSON LoggerType = "json"
)

func NewSlog(t LoggerType) {
	var logger *slog.Logger

	levelOpts := map[Level]slog.Level{
		Debug: slog.LevelDebug,
		Info:  slog.LevelInfo,
		Warn:  slog.LevelWarn,
		Error: slog.LevelError,
	}

	var handler slog.Handler

	switch t {
	case JSON:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: levelOpts[GetLevel()]})
		break
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: levelOpts[GetLevel()]})
	}

	logger = slog.New(handler)

	slog.SetDefault(logger)
}

func GetLevel() Level {
	if level, exists := os.LookupEnv("LOG_LEVEL"); exists {
		switch Level(level) {
		case Debug:
			return Debug
		case Info:
			return Info
		case Warn:
			return Warn
		case Error:
			return Error
		default:
			return Info
		}
	}
	return Info
}
