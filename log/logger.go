package log

import (
    "fmt"
    "sync"
)

type InfoLogger interface {
    Info(...interface{}) error
}

type UnsafeInfoLogger interface {
    Info(...interface{})
}

type ErrorLogger interface {
    Error(...interface{}) error
}

type UnsafeErrorLogger interface {
    Error(...interface{})
}

type InfoErrorLogger interface {
    InfoLogger
    ErrorLogger
}
type UnsafeInfoErrorLogger interface {
    UnsafeInfoLogger
    UnsafeErrorLogger
}

type ignoringInfoError struct {
    ielogger InfoErrorLogger
}

func Unsafe(ie InfoErrorLogger) *ignoringInfoError {
    return &ignoringInfoError{
        ielogger: ie,
    }
}

func (l *ignoringInfoError) Info(args ...interface{}) {
    _ = l.ielogger.Info(args...)
}

func (l *ignoringInfoError) Error(args ...interface{}) {
    _ = l.ielogger.Error(args...)
}

type MultiInfoError struct {
    loggers []InfoErrorLogger
}

func Multi(loggers ...InfoErrorLogger) *MultiInfoError {
    return &MultiInfoError{
        loggers: loggers,
    }
}

func (l *MultiInfoError) Info(args ...interface{}) error {
    return l.log(func(logger InfoErrorLogger) func(...interface{}) error {
        return logger.Info
    }, args...)
}

func (l *MultiInfoError) Error(args ...interface{}) error {
    return l.log(func(logger InfoErrorLogger) func(...interface{}) error {
        return logger.Error
    }, args...)
}

type multipleErrors []error

func newMultipleErrors(errs ...error) multipleErrors {
    return errs
}

func (e multipleErrors) Error() string {
    return fmt.Sprintf("multiple errors: %v", e)
}

func (l *MultiInfoError) log(f func(InfoErrorLogger) func(...interface{}) error, args ...interface{}) error {
    var (
        wg   sync.WaitGroup
        errs []error
    )
    wg.Add(len(l.loggers))
    for _, logger := range l.loggers {
        go func(lg InfoErrorLogger) {
            defer wg.Done()
            if err := f(lg)(args...); err != nil {
                errs = append(errs, err)
            }
        }(logger)
    }
    wg.Wait()
    if errs != nil {
        return newMultipleErrors(errs...)
    }
    return nil
}

type PropagateInfoError struct {
    main     InfoErrorLogger
    insurers []ErrorLogger
}

func NewPropagateInfoError(main InfoErrorLogger, insurers ...ErrorLogger) *PropagateInfoError {
    return &PropagateInfoError{
        main:     main,
        insurers: insurers,
    }
}

func (l *PropagateInfoError) Info(args ...interface{}) error {
    return l.log(func(logger InfoErrorLogger) func(...interface{}) error {
        return logger.Info
    }, args...)
}

func (l *PropagateInfoError) Error(args ...interface{}) error {
    return l.log(func(logger InfoErrorLogger) func(...interface{}) error {
        return logger.Error
    }, args...)
}

func (l *PropagateInfoError) log(f func(InfoErrorLogger) func(...interface{}) error, args ...interface{}) error {
    err := f(l.main)(args...)
    lastErr := err
    for i := 0; i < len(l.insurers) && lastErr != nil; i++ {
        lastErr = l.insurers[i].Error(err)
    }
    return err
}
