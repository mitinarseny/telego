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
    usersCollectionName = "users"
)

type UsersRepo struct {
    this *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) (*UsersRepo, error) {
    r := &UsersRepo{
        this: db.Collection(usersCollectionName),
    }
    if _, err := r.this.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{"id", 1}},
        Options: options.Index().SetUnique(true),
    }); err != nil {
        return nil, errors.Wrapf(err, "unable to create index on %s", usersCollectionName)
    }
    return r, nil
}

func (r *UsersRepo) Create(ctx context.Context, users ...*tglog.User) ([]*tglog.User, error) {
    models := make([]interface{}, 0, len(users))
    for _, u := range users {
        models = append(models, u)
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return users, nil
}

func (r *UsersRepo) CreateIfNotExist(ctx context.Context, users ...*tglog.User) ([]*tglog.User, error) {
    models := make([]mongo.WriteModel, 0, len(users))
    for _, user := range users {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"id", user.ID},
            }).SetUpdate(user).SetUpsert(true))
    }
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return users, nil
}
