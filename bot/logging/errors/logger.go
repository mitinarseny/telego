package errors

import "sync"

type Logger interface{
    Log(error)
}

type CompositeErrorLogger struct {
    loggers []Logger
}

func NewCompositeErrorLogger(loggers ...Logger) *CompositeErrorLogger {
    return &CompositeErrorLogger{
        loggers: loggers,
    }
}

func (l *CompositeErrorLogger) Log(err error) {
    var wg sync.WaitGroup
    wg.Add(len(l.loggers))
    for _, logger := range l.loggers {
        go func() {
            defer wg.Done()
            logger.Log(err)
        }()
    }
    wg.Wait()
}
