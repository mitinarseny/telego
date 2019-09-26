package mongo

import (
    errors2 "github.com/mitinarseny/telego/bot/logging/errors"
    "github.com/pkg/errors"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    errLogCollectionName = "errors"
)

type Logger struct {
    this  *mongo.Collection
    spare errors2.Logger
}

func NewErrorLogger(db *mongo.Database, spareLogger errors2.Logger) *Logger {
    return &Logger{
        this:  db.Collection(errLogCollectionName),
        spare: spareLogger,
    }
}

func (l *Logger) Log(err error) {
    if err := l.log(err); err != nil {
        l.spare.Log(err)
    }
}

func (l *Logger) log(err error) error {
    return errors.New("mongo logger is not implemented yet")
}
