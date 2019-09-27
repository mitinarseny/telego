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
    return l.log("error", args...)
}

func (l *Logger) Info(args ...interface{}) error {
    return l.log("info", args...)
}

func (l *Logger) log(typ string, args ...interface{}) error {
    _, err := l.this.InsertOne(context.Background(), bson.D{
        {"created_at", time.Now()},
        {"type", typ},
        {"data", args},
    })
    return err
}
