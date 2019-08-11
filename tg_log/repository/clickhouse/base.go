package clickhouse

import (
    "context"
    "database/sql"
)

const (
    logDB = "log"
    updatesTable = "updates"
)

type baseRepository struct {
    conn      *sql.DB
}

func newBaseRepository(db *sql.DB) *baseRepository {
    return &baseRepository{
        conn:      db,
    }
}

type txFunc func(tx *sql.Tx) error

func (repo *baseRepository) makeTransaction(ctx context.Context, f txFunc) (err error) {
    tx, err := repo.conn.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p) // re-throw panic after Rollback
        } else if err != nil {
            _ = tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()
    err = f(tx)
    return err
}
