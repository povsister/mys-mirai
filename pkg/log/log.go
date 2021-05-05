package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"time"
)

var defaultLogger zerolog.Logger

var (
	Level = zerolog.InfoLevel
)

func InitLogger() {
	currentDir := filepath.Dir(os.Args[0])
	_ = os.MkdirAll(filepath.Join(currentDir, "logs"), 0755)

	fd, err := NewRotator(filepath.Join(currentDir, "logs", "debug.log"))
	if err != nil {
		stdErr("WARNING: failed to create debug log file %v", err)
	}

	lw := &levelWriter{
		ConsoleWriter: zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true, TimeFormat: time.RFC3339}, l: Level,
	}
	var w io.Writer = lw
	if fd != nil {
		w = zerolog.MultiLevelWriter(lw, fd)
	}

	defaultLogger = zerolog.New(w).With().Timestamp().Logger().Level(zerolog.TraceLevel)
}

type levelWriter struct {
	zerolog.ConsoleWriter
	l zerolog.Level
}

func (w *levelWriter) WriteLevel(l zerolog.Level, p []byte) (int, error) {
	if l >= w.l {
		return w.Write(p)
	}
	return len(p), nil
}

func Trace() *zerolog.Event {
	return defaultLogger.Trace()
}

func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

func Info() *zerolog.Event {
	return defaultLogger.Info()
}

func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

func Error() *zerolog.Event {
	return defaultLogger.Error()
}

func Fatal() *zerolog.Event {
	return defaultLogger.Fatal()
}
