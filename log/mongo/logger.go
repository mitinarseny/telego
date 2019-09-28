package mongo

import (
    "context"
    "fmt"

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
    var formattedArgs []interface{}
    for _, a := range args {
        ca := a
        switch at := a.(type) {
        case error:
            ca = at.Error()
        case fmt.Stringer:
            ca = at.String()
        }
        formattedArgs = append(formattedArgs, ca)
    }
    _, err := l.this.InsertOne(context.Background(), bson.D{
        {"level", lvl},
        {"data", formattedArgs},
    })
    return err
}
