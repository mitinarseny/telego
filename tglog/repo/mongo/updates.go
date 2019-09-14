package mongo

import (
    "context"
    "time"

    "github.com/mitinarseny/telego/tglog/repo"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    updatesCollectionName = "updates"
)

type UpdatesRepository struct {
    collection *mongo.Collection
}

func NewUpdatesRepository(db *mongo.Database) *UpdatesRepository {
    return &UpdatesRepository{
        collection: db.Collection(updatesCollectionName),
    }
}

func (repo *UpdatesRepository) Create(ctx context.Context, updates []repo.Update) ([]repo.Update, error) {
    upds := make([]interface{}, 0, len(updates))
    for _, update := range updates {
        ca := time.Now()
        update.CreatedAt = &ca
        upds = append(upds, update)
    }
    _, err := repo.collection.InsertMany(ctx, upds)
    if err != nil {
        return nil, err
    }
    return updates, nil
}
