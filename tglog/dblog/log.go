package dblog

import (
    "context"

    "github.com/mitinarseny/telego/tglog/repo"
)

type RepoLogger struct {
    UpdatesRepo repo.UpdatesRepo
}

func (l *RepoLogger) LogUpdates(updates ...*repo.Update) error {
    _, err := l.UpdatesRepo.Create(context.Background(), updates...)
    return err
}
