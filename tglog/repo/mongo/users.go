package mongo

import (
    "context"

    "github.com/mitinarseny/telego/tglog/repo"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/mongo"
)

const (
    usersCollectionName = "users"
)

type UsersRepo struct {
    collection *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
    return &UsersRepo{
        collection: db.Collection(usersCollectionName),
    }
}

func (r *UsersRepo) Create(ctx context.Context, users ...*repo.User) ([]*repo.User, error) {
    usrs := make([]interface{}, 0, len(users))
    for _, user := range users {
        usrs = append(usrs, user)
    }
    _, err := r.collection.InsertMany(ctx, usrs)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "UsersRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "UsersRepo",
        "status":  "CREATED",
        "count":   len(users),
    }).Info()
    return users, nil
}
