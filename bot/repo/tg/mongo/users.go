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

func (r *UsersRepo) Create(ctx context.Context, users ...*tg.User) ([]*tg.User, error) {
    models := make([]interface{}, 0, len(users))
    for _, u := range users {
        u.BaseModel.CreatedAt = time.Now()
        models = append(models, u)
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return users, nil
}

func (r *UsersRepo) CreateIfNotExist(ctx context.Context, users ...*tg.User) ([]*tg.User, error) {
    models := make([]mongo.WriteModel, 0, len(users))
    for _, user := range users {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", user.ID},
            }).SetUpdate(user).SetUpsert(true))
    }
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return users, nil
}
