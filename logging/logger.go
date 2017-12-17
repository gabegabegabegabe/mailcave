package logging

import (
	"fmt"
	"path"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger writes logs to a file.
type Logger struct {
	logToStdOut bool
	*lumberjack.Logger
}

// NewLogger creates a Logger.
func NewLogger(appName string, logDir string, maxSizeBytes int, maxBackups int, maxAgeDays int, logToStdOut bool) *Logger {
	return &Logger{
		Logger: &lumberjack.Logger{
			Filename:   path.Join(logDir, appName+".log"),
			MaxSize:    maxSizeBytes,
			MaxBackups: maxBackups,
			MaxAge:     maxAgeDays,
			LocalTime:  true,
		},
		logToStdOut: logToStdOut,
	}
}

// Printf allows formatted string logging.
func (l *Logger) Printf(format string, a ...interface{}) (n int, err error) {
	t := time.Now()
	zone, _ := t.Zone()
	timeStr := t.Format("Mon Jan 2 15:04:05 " + zone + " 2006")

	msg := fmt.Sprintf(timeStr+"  "+format+"\n", a...)

	if l.logToStdOut {
		fmt.Printf(msg)
	}

	return l.Write([]byte(msg))
}
