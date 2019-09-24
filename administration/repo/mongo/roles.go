package mongo

import (
    "context"
    "errors"
    "fmt"

    "github.com/mitinarseny/telego/administration/repo"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    rolesCollectionName = "roles"
)

type RolesRepo struct {
    this *mongo.Collection
}

func (r *RolesRepo) Create(ctx context.Context, roles ...*repo.Role) ([]*repo.Role, error) {
    rls := make([]interface{}, 0, len(roles))
    for _, role := range roles {
        rls = append(rls, role)
    }
    res, err := r.this.InsertMany(ctx, rls)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "AdminsRepo",
        "status":  "CREATED",
        "count":   len(res.InsertedIDs),
    }).Info()
    return roles, nil
}

func (r *RolesRepo) CreateIfNotExists(ctx context.Context, roles ...*repo.Role) ([]*repo.Role, error) {
    models := make([]mongo.WriteModel, 0, len(roles))
    for _, role := range roles {
        models = append(models,
            mongo.NewUpdateOneModel().
                SetFilter(bson.D{
                    {"_id", role.Name},
                }).
                SetUpdate(role).
                SetUpsert(true))
    }
    res, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "CreateIfNotExists",
        }).Error(err)
        return nil, err
    }
    if res.UpsertedCount > 0 {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "status":  "CREATED",
            "count":   res.UpsertedCount,
        }).Info()
    }
    return roles, nil
}

func (r *RolesRepo) GetAll(ctx context.Context) ([]*repo.Role, error) {
    cursor, err := r.this.Find(ctx, bson.D{})
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "GetAll",
        }).Error(err)
        return nil, err
    }
    var roles []*repo.Role
    if err := cursor.All(ctx, &roles); err != nil {
        return nil, err
    }
    return roles, nil
}

func (r *RolesRepo) GetByName(ctx context.Context, name string) (*repo.Role, error) {
    roles, err := r.GetByNames(ctx, name)
    if err != nil {
        return nil, err
    }
    switch {
    case len(roles) == 0:
        return nil, errors.New(fmt.Sprintf("role %q not found", name))
    case len(roles) > 1:
        return nil, errors.New(fmt.Sprintf("more than one role with name %q found", name))
    }
    return roles[0], nil
}

func (r *RolesRepo) GetByNames(ctx context.Context, names ...string) ([]*repo.Role, error) {
    cursor, err := r.this.Find(ctx, bson.D{
        {"_id", bson.D{
            {"$in", names},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "GetByNames",
        }).Error(err)
        return nil, err
    }
    roles := make([]*repo.Role, 0, len(names))
    if err := cursor.All(ctx, &roles); err != nil {
        return nil, err
    }
    return roles, nil
}

func (r *RolesRepo) AddScopes(ctx context.Context, scopes []repo.Scope, names ...string) ([]*repo.Role, error) {
    res, err := r.this.UpdateMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", names},
        }},
    }, bson.D{
        {"$addToSet", bson.D{
            {"scopes", bson.D{
                {"$each", scopes},
            }},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "AddScopes",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "RolesRepo",
        "action":  "AddScopes",
        "count":   res.ModifiedCount,
    }).Info()
    return r.GetByNames(ctx, names...)
}

func (r *RolesRepo) SetScopes(ctx context.Context, scopes []repo.Scope, names ...string) ([]*repo.Role, error) {
    res, err := r.this.UpdateMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", names},
        }},
    }, bson.D{
        {"$set", bson.D{
            {"scopes", scopes},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "AddScopes",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "RolesRepo",
        "action":  "AddScopes",
        "count":   res.ModifiedCount,
    }).Info()
    return r.GetByNames(ctx, names...)
}

func (r *RolesRepo) DeleteByNames(ctx context.Context, names ...string) error {
    res, err := r.this.DeleteMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", names},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "RolesRepo",
            "action":  "DeleteByNames",
        }).Error(err)
        return err
    }
    log.WithFields(log.Fields{
        "context": "RolesRepo",
        "action":  "DeleteByNames",
        "count":   res.DeletedCount,
    }).Info()
    return nil
}

func NewRolesRepo(db *mongo.Database) *RolesRepo {
    return &RolesRepo{
        this: db.Collection(rolesCollectionName),
    }
}
