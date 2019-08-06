package ch_log

import (
    "database/sql"
    "errors"
    "fmt"
    "net/url"
    "time"

    "github.com/charithe/timedbuf"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    log "github.com/sirupsen/logrus"
)

const (
    buffSize = 1000 // for buffer table
    buffTime = 10 * time.Second
)

type UpdatesLogger struct {
    tb   *timedbuf.TimedBuf
    conn *sql.DB
}

func NewUpdatesLogger(host string, port uint, username, password, dbName, tableName string) (*UpdatesLogger, error) {
    chURL := url.URL{
        Scheme: "tcp",
        User:   url.UserPassword(username, password),
        Host:   fmt.Sprintf("%s:%d", host, port),
    }
    conn, err := sql.Open("clickhouse", chURL.String())
    if err != nil {
        return nil, err
    }

    if err := conn.Ping(); err != nil {
        return nil, err
    }





    _, err = conn.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %[1]s.%[1]s_buffer AS %[1]s.%[1]s engine=Buffer(
            %[1]s, %[2]s, 16, 10, 100, 1000, 100000, 10000000, 100000000
        )`, dbName, tableName))
    if err != nil {
        return nil, err
    }

    return &UpdatesLogger{
        tb: timedbuf.New(buffSize, buffTime, func(items []interface{}) {
            switch updates := interface{}(items).(type) {
            case []*tgbotapi.Update:
                if err := flush(conn, updates); err != nil {
                    log.WithFields(log.Fields{
                        "context": "LOG",
                        "driver":  "clickhouse",
                    }).Error(err)
                    return
                }
            default:
                log.WithFields(log.Fields{
                    "context": "LOG",
                    "driver":  "clickhouse",
                }).Error(errors.New("unable to log not []*tgbotapi.Update"))
                return
            }
        }),
        conn: conn,
    }, nil
}

func (ul *UpdatesLogger) LogUpdates(updates tgbotapi.UpdatesChannel) {
    for update := range updates {
        ul.tb.Put(update)
    }
}

func flush(conn *sql.DB, updates []*tgbotapi.Update) error {
    log.WithFields(log.Fields{
        "context":      "LOG",
        "driver":       "clickhouse",
        "updatesCount": len(updates),
    }).Info("Flushing updates log to ClickHouse")

    tx, err := conn.Begin()
    if err != nil {
        return err
    }
    query := "INSERT INTO () VALUES "
    for update := range updates {
        query
    }

    stmt, err := tx.Prepare("INSERT INTO () VALUES ()") // TODO: tables, values, buffer table
    stmt.
    if err != nil {
        return err
    }
    defer func() {
        _ = stmt.Close()
    }()

    for update := range updates {
        if _, err = stmt.Exec(); err != nil {
            return err
        }
    }
    if err := tx.Commit(); err != nil {
        return err
    }
    return nil
}
