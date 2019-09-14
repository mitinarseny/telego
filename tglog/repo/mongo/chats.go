package mongo

import (
    "context"

    "github.com/mitinarseny/telego/tglog/repo"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/mongo"
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
