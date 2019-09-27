package stderr

import (
    log "github.com/sirupsen/logrus"
)

type Logger struct {
    entry *log.Entry
}

func NewErrorLogger(l *log.Logger) *Logger {
    return &Logger{
        entry: l.WithFields(log.Fields{
            "context": "BOT",
        }),
    }
}

func (l *Logger) Error(args ...interface{}) error {
    l.entry.Error(args...)
    return nil
}

func (l *Logger) Info(args ...interface{}) error {
    l.entry.Info(args...)
    return nil
}
