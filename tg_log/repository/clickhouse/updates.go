package clickhouse

import (
    "context"
    "database/sql"
    "fmt"

    _ "github.com/kshvakov/clickhouse"
    "github.com/mitinarseny/telego/tg_log/repository"
)

const (
    updatesUpdateIDColumn = "update_id"
)

type UpdatesRepository struct {
    *baseRepository
}

func NewUpdatesRepository(db *sql.DB) *UpdatesRepository {
    return &UpdatesRepository{
        baseRepository: newBaseRepository(db),
    }
}

func (repo *UpdatesRepository) Create(ctx context.Context, updates ...*repository.Update) error {
    return repo.makeTransaction(ctx, func(tx *sql.Tx) error {
        stmt, err := tx.Prepare(fmt.Sprintf(
            "INSERT INTO %s.%s (%s) VALUES (?)",
            logDB,
            updatesTable,
            updatesUpdateIDColumn,
        ))
        if err != nil {
            return err
        }
        var count int
        for _, u := range updates {
            if _, err := stmt.Exec(
                u.UpdateID,
            ); err != nil {
                return err
            }
            count++
        }
        return nil
    })
}
