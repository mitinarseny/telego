package dblog

import (
    "context"

    "github.com/mitinarseny/telego/tglog/repo"
)

type DBLogger struct {
    UpdatesRepo repo.UpdatesRepo
}

func (l *DBLogger) LogUpdates(updates []repo.Update) error {
    _, err := l.UpdatesRepo.Create(context.Background(), updates)
    return err
}
