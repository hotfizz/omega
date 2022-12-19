package alt

import (
	"fmt"
	"io"
	"time"
)

type LogLevel int

func (c LogLevel) String() string {
	switch c {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "DEBUG"
	}
}

const (
	LevelDebug LogLevel = iota
	LevelInfo           = 1
	LevelWarn           = 2
	LevelError          = 3
)
const (
	logMsg = "[%s] %s\n"
)

type stdLogger struct {
	level  LogLevel
	writer io.Writer
}

func NewStdLogger(level LogLevel, writer io.Writer) Logger {
	return &stdLogger{
		level:  level,
		writer: writer,
	}
}

func (c *stdLogger) Debug(args ...interface{}) {
	if LevelDebug < c.level {
		return
	}
	msg := fmt.Sprintf(logMsg, c.level.String(), time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprint(c.writer, append([]interface{}{msg}, args...)...)
}

func (c stdLogger) Info(args ...interface{}) {
	if LevelInfo < c.level {
		return
	}
	msg := fmt.Sprintf(logMsg, c.level.String(), time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprint(c.writer, append([]interface{}{msg}, args...)...)
}

func (c stdLogger) Warn(args ...interface{}) {
	if LevelWarn < c.level {
		return
	}
	msg := fmt.Sprintf(logMsg, c.level.String(), time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprint(c.writer, append([]interface{}{msg}, args...)...)
}

func (c stdLogger) Error(args ...interface{}) {
	if LevelError < c.level {
		return
	}
	msg := fmt.Sprintf(logMsg, c.level.String(), time.Now().Format(time.RFC3339))
	_, _ = fmt.Fprint(c.writer, append([]interface{}{msg}, args...)...)
}
