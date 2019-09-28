package mongo

import (
    "context"

    "github.com/mitinarseny/telego/bot/tglog"
    "github.com/pkg/errors"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    chatsCollectionName = "chats"
)

type ChatsRepo struct {
    this *mongo.Collection
}

func NewChatsRepo(db *mongo.Database) (*ChatsRepo, error) {
    r := &ChatsRepo{
        this: db.Collection(chatsCollectionName),
    }
    if _, err := r.this.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{"id", 1}},
        Options: options.Index().SetUnique(true),
    }); err != nil {
        return nil, errors.Wrapf(err, "unable to create index on %s", chatsCollectionName)
    }
    return r, nil
}

func (r *ChatsRepo) Create(ctx context.Context, chats ...*tglog.Chat) ([]*tglog.Chat, error) {
    models := make([]interface{}, 0, len(chats))
    for _, c := range chats {
        models = append(models, c)
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return chats, nil
}

func (r *ChatsRepo) CreateIfNotExist(ctx context.Context, chats ...*tglog.Chat) ([]*tglog.Chat, error) {
    models := make([]mongo.WriteModel, 0, len(chats))
    for _, c := range chats {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"id", c.ID},
            }).SetUpdate(c).SetUpsert(true))
    }
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return chats, nil
}
