package mongo

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    errLogCollectionName = "logs"
)

type Logger struct {
    this *mongo.Collection
}

func NewErrorLogger(db *mongo.Database) *Logger {
    return &Logger{
        this: db.Collection(errLogCollectionName),
    }
}

func (l *Logger) Error(args ...interface{}) error {
    return l.log(errorLevel, args...)
}

func (l *Logger) Info(args ...interface{}) error {
    return l.log(infoLevel, args...)
}

const (
    errorLevel level = "error"
    infoLevel  level = "info"
)

type level string

func (l *Logger) log(lvl level, args ...interface{}) error {
    _, err := l.this.InsertOne(context.Background(), bson.D{
        {"created_at", time.Now()},
        {"level", lvl},
        {"data", args},
    })
    return err
}
