package clickhouse

import (
    "database/sql"
    "net/url"
    "strconv"
    "time"

    tb "github.com/charithe/timedbuf"
    "github.com/kshvakov/clickhouse"
    "github.com/mitinarseny/telego/tg_log/tg_types"
    log "github.com/sirupsen/logrus"
)

const (
    buffSize     = 10000
    buffDuration = 5 * time.Second // TODO: minute?
)

var (
    logEntry = log.WithFields(log.Fields{
        "context": "ClickHouseUpdatesLogger",
    })
)

type BufferedUpdateLogger struct {
    tb *tb.TimedBuf
}

func NewBufferedUpdateLogger(host string, port int, username, password, dbName string) (*BufferedUpdateLogger, error) {
    connURL := url.URL{
        Scheme: "tcp",
        Host:   host,
        Path:   dbName,
    }
    if port != 0 {
        connURL.Host += ":" + strconv.Itoa(port)
    }
    if username != "" {
        connURL.RawQuery += "&username=" + username
    }
    if password != "" {
        connURL.RawQuery += "&password=" + password
    }

    db, err := sql.Open("clickhouse", connURL.String())
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        if exception, ok := err.(*clickhouse.Exception); ok {
            logEntry.Error(exception) // TODO: rm this shit
        }
        return nil, err
    }

    return &BufferedUpdateLogger{
        tb: tb.New(buffSize, buffDuration, func(items []interface{}) {
            tx, err := db.Begin()
            if err != nil {
                logEntry.WithFields(log.Fields{
                    "action": "BEGIN TRANSACTION",
                }).Error(err)
                return
            }
            stmt, err := tx.Prepare("INSERT INTO log.updates (update_id) VALUES (?)")
            if err != nil {
                logEntry.WithFields(log.Fields{
                    "action": "PREPARE STATEMENT",
                }).Error(err)
                return
            }
            for _, i := range items {
                u, ok := i.(*tg_types.Update)
                if !ok {
                    logEntry.WithFields(log.Fields{
                        "count": len(items),
                        "what":  i,
                    }).Error("Unable to flush not Update")
                    continue
                }
                _, err := stmt.Exec(
                    u.UpdateID,
                )
                if err != nil {
                    log.WithFields(log.Fields{
                        "update": *u,
                        "action": "EXEC",
                    }).Error(err)
                    continue
                }
            }

            if err := tx.Commit(); err != nil {
                logEntry.WithFields(log.Fields{
                    "count":  len(items),
                    "action": "INSERT",
                }).Error(err)
                if err := tx.Rollback(); err != nil {
                    logEntry.WithFields(log.Fields{
                        "action": "ROLLBACK",
                    }).Error(err)
                    return
                }
                logEntry.WithFields(log.Fields{
                    "action": "ROLLBACK",
                }).Info()
                return
            }
            logEntry.WithFields(log.Fields{
                "action": "FLUSH",
                "status": "SUCCESS",
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
