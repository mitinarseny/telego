package clickhouse

import (
    "database/sql"
    "fmt"
    "net/url"
    "time"

    tb "github.com/charithe/timedbuf"
    "github.com/kshvakov/clickhouse"
    "github.com/mitinarseny/telego/ulog/tg_types"
    log "github.com/sirupsen/logrus"
)

const (
    buffSize     = 10000
    buffDuration = 1 * time.Minute
)

type BufferedUpdateLogger struct {
    tb *tb.TimedBuf
}

func NewBufferedUpdateLogger(host string, port int, username, password string) (*BufferedUpdateLogger, error) {
    connURL := url.URL{
        Scheme: "tcp",
        Host:   host,
    }
    if port != 0 {
        connURL.Host += fmt.Sprintf(":%d", port)
    }
    if username != "" {
        if password != "" {
            connURL.User = url.UserPassword(username, password)
        } else {
            connURL.User = url.User(username)
        }
    }
    db, err := sql.Open("clickhouse", connURL.String())
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        if exception, ok := err.(*clickhouse.Exception); ok {
            log.Error(exception)
        }
        return nil, err
    }

    if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS updates"); err != nil {
        return nil, err
    }

    return &BufferedUpdateLogger{
        tb: tb.New(buffSize, buffDuration, func(items []interface{}) {
            log.WithFields(log.Fields{
                "context": "ClickHouseUpdatesLogger",
                "count":   len(items),
            }).Info()
        }),
    }, nil
}

func (l *BufferedUpdateLogger) LogUpdate(u *tg_types.Update) error {
    l.tb.Put(u)
    return nil
}

func (l *BufferedUpdateLogger) Close() {
    l.tb.Close()
}
