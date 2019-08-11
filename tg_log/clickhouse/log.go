package clickhouse

import (
    "context"
    "database/sql"
    "time"

    tb "github.com/charithe/timedbuf"
    _ "github.com/kshvakov/clickhouse"
    "github.com/mitinarseny/telego/tg_log/repository"
    "github.com/mitinarseny/telego/tg_log/repository/clickhouse"
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

type BufferedUpdatesLogger struct {
    tb *tb.TimedBuf
    r updatesLogRepository
}

func NewBufferedUpdateLogger(db *sql.DB) (*BufferedUpdatesLogger, error) {
    chRepo := clickhouse.NewUpdatesRepository(db)
    return &BufferedUpdatesLogger{
        r: chRepo,
        tb: tb.New(buffSize, buffDuration, func(items []interface{}) {
            updates := make([]*repository.Update, 0, len(items))
            for _, i := range items {
                u, ok := i.(*repository.Update)
                if !ok {
                    logEntry.WithFields(log.Fields{
                        "what":  i,
                    }).Error("Unable to flush not Update")
                    return
                }
                updates = append(updates, u)
            }
            if err := chRepo.Create(context.Background(), updates...); err != nil {
                logEntry.WithFields(log.Fields{
                    "action": "FLUSH",
                    "status": "ERROR",
                }).Error(err)
                return
            }
            logEntry.WithFields(log.Fields{
                "action": "FLUSH",
                "status": "SUCCESS",
                "count":  len(items),
            }).Info()
        }),
    }, nil
}

func (l *BufferedUpdatesLogger) LogUpdate(u *repository.Update) error {
    l.tb.Put(u)
    return nil
}

func (l *BufferedUpdatesLogger) Close() {
    l.tb.Close()
}

type updatesLogRepository interface {

}
