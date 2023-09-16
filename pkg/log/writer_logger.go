package log

import (
	"io"
)

type writerLoggerAdapter struct {
	logger Logger
	level  LogLevel
}

func NewWriterLoggerAdapter(logger Logger, level LogLevel) io.Writer {
	return &writerLoggerAdapter{
		logger,
		level,
	}
}

func (w *writerLoggerAdapter) Write(p []byte) (n int, err error) {
	switch w.level {
	case DEBUG:
		w.logger.Debug(string(p))
	case INFO:
		w.logger.Info(string(p))
	case WARNING:
		w.logger.Warn(string(p))
	default:
		w.logger.Debug(string(p))
	}
	return len(p), nil
}
