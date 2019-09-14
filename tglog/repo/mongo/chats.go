package mongo

import (
    "context"

    "github.com/mitinarseny/telego/tglog/repo"
    log "github.com/sirupsen/logrus"
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

func (r *ChatsRepo) Create(ctx context.Context, chats ...*repo.Chat) ([]*repo.Chat, error) {
    chts := make([]interface{}, 0, len(chats))
    for _, chat := range chats {
        chts = append(chts, chat)
    }
    _, err := r.this.InsertMany(ctx, chts)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "ChatsRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "ChatsRepo",
        "status":  "CREATED",
        "count":   len(chats),
    }).Info()
    return chats, nil
}

func (r *ChatsRepo) CreateIfNotExists(ctx context.Context, chats ...*repo.Chat) ([]*repo.Chat, error) {
    models := make([]mongo.WriteModel, 0, len(chats))
    for _, chat := range chats {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", chat.ID},
            }).SetUpdate(chat).SetUpsert(true))
    }
    res, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        log.WithFields(log.Fields{
            "context": "ChatsRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    if res.UpsertedCount > 0 {
        log.WithFields(log.Fields{
            "context": "ChatsRepo",
            "status":  "CREATED",
            "count":   res.UpsertedCount,
        }).Info()
    }
    return chats, nil
}
