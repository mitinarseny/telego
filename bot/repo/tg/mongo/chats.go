package mongo

import (
    "context"
    "time"

    "github.com/mitinarseny/telego/bot/repo/tg"
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

func NewChatsRepo(db *mongo.Database) *ChatsRepo {
    return &ChatsRepo{
        this: db.Collection(chatsCollectionName),
    }
}

func (r *ChatsRepo) Create(ctx context.Context, chats ...*tg.Chat) ([]*tg.Chat, error) {
    models := make([]interface{}, 0, len(chats))
    for _, c := range chats {
        c.BaseModel.CreatedAt = time.Now()
        models = append(models, c)
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return chats, nil
}

func (r *ChatsRepo) CreateIfNotExist(ctx context.Context, chats ...*tg.Chat) ([]*tg.Chat, error) {
    models := make([]mongo.WriteModel, 0, len(chats))
    for _, c := range chats {
        c.BaseModel.CreatedAt = time.Now()
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", c.ID},
            }).SetUpdate(c).SetUpsert(true))
    }
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return chats, nil
}
