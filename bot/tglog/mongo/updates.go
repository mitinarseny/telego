package mongo

import (
    "context"
    "time"

    "github.com/mitinarseny/telego/bot/tglog"
    "github.com/pkg/errors"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    updatesCollectionName = "updates"
)

type UpdatesRepo struct {
    this  *mongo.Collection
    users tglog.UsersRepo
    chats tglog.ChatsRepo
}

func NewUpdatesRepo(db *mongo.Database, ur tglog.UsersRepo, cr tglog.ChatsRepo) *UpdatesRepo {
    return &UpdatesRepo{
        this:  db.Collection(updatesCollectionName),
        users: ur,
        chats: cr,
    }
}

func (r *UpdatesRepo) Create(ctx context.Context, updates ...*tglog.Update) ([]*tglog.Update, error) {
    models := make([]interface{}, 0, len(updates))
    users := make([]*tglog.User, 0, len(updates))
    chats := make([]*tglog.Chat, 0, len(updates))
    for _, u := range updates {
        if from := u.From(); from != nil {
            users = append(users, from)
        }
        if chat := u.Chat(); chat != nil {
            chats = append(chats, chat)
        }
        u.BaseModel.CreatedAt = time.Now()
        models = append(models, u)
    }
    if _, err := r.users.CreateIfNotExist(ctx, users...); err != nil {
        return nil, err
    }
    if _, err := r.chats.CreateIfNotExist(ctx, chats...); err != nil {
        return nil, err
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return updates, nil
}

func (r *UpdatesRepo) CreateIfNotExist(ctx context.Context, updates ...*tglog.Update) ([]*tglog.Update, error) {
    return nil, errors.New("UpdatesRepo.CreateIfNotExist is not implemented yet")
}
