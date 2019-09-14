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
    usersCollectionName = "users"
)

type UsersRepo struct {
    this *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
    return &UsersRepo{
        this: db.Collection(usersCollectionName),
    }
}

func (r *UsersRepo) Create(ctx context.Context, users ...*repo.User) ([]*repo.User, error) {
    usrs := make([]interface{}, 0, len(users))
    for _, user := range users {
        usrs = append(usrs, user)
    }
    _, err := r.this.InsertMany(ctx, usrs)
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

func (r *UsersRepo) CreateIfNotExists(ctx context.Context, users ...*repo.User) ([]*repo.User, error) {
    models := make([]mongo.WriteModel, 0, len(users))
    for _, user := range users {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", user.ID},
            }).SetUpdate(user).SetUpsert(true))
    }
    res, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        log.WithFields(log.Fields{
            "context": "UsersRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    if res.UpsertedCount > 0 {
        log.WithFields(log.Fields{
            "context": "UsersRepo",
            "status":  "CREATED",
            "count":   res.UpsertedCount,
        }).Info()
    }
    return users, nil
}
