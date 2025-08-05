package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Config struct {
	Level  string
	Output string
}

func Init(cfg Config) error {
	// Уровень
	level, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	Log.SetLevel(level)

	// Вывод
	switch cfg.Output {
	case "stdout":
		Log.SetOutput(os.Stdout)
	case "stderr":
		Log.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		Log.SetOutput(io.MultiWriter(f, os.Stdout))
	}

	// Формат JSON
	Log.SetFormatter(&logrus.JSONFormatter{})

	return nil
}
