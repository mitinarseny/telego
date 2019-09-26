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

func (l *Logger) Log(err error) {
    l.entry.Error(err)
}
