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
    adminsCollectionName = "admins"
)

type AdminsRepo struct {
    this  *mongo.Collection
    roles repo.RolesRepo
}

type AdminsRepoDependencies struct {
    Roles repo.RolesRepo
}

func (d *AdminsRepoDependencies) Validate() error {
    shouldBeNotNil := [...]interface{}{
        d.Roles,
    }
    for _, e := range shouldBeNotNil {
        if e == nil {
            return errors.New(fmt.Sprintf("%T should be not nil", e))
        }
    }
    return nil
}

func NewAdminsRepo(db *mongo.Database, deps *AdminsRepoDependencies) (*AdminsRepo, error) {
    if err := deps.Validate(); err != nil {
        return nil, err
    }
    return &AdminsRepo{
        this:  db.Collection(adminsCollectionName),
        roles: deps.Roles,
    }, nil
}

func (r *AdminsRepo) Create(ctx context.Context, admins ...*repo.Admin) ([]*repo.Admin, error) {
    adms := make([]interface{}, 0, len(admins))
    for _, admin := range admins {
        adms = append(adms, admin)
    }
    res, err := r.this.InsertMany(ctx, adms)
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "CREATE",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "AdminsRepo",
        "status":  "CREATED",
        "count":   len(res.InsertedIDs),
    }).Info()
    return admins, nil
}

func (r *AdminsRepo) CreateIfNotExists(ctx context.Context, admins ...*repo.Admin) ([]*repo.Admin, error) {
    models := make([]mongo.WriteModel, 0, len(admins))
    for _, admin := range admins {
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", admin.ID},
            }).SetUpdate(admin).SetUpsert(true))
    }
    res, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "CreateIfNotExists",
        }).Error(err)
        return nil, err
    }
    if res.UpsertedCount > 0 {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "status":  "CREATED",
            "count":   res.UpsertedCount,
        }).Info()
    }
    return admins, nil
}

func (r *AdminsRepo) ChangeRoleByIDs(ctx context.Context, role *repo.Role, adminIDs ...int64) ([]*repo.Admin, error) {
    res, err := r.this.UpdateMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    }, bson.D{
        {"$set", bson.D{
            {"role", role.Name},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "ChangeRoleByIDs",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "UpdatesRepo",
        "action":  "ChangeRoleByIDs",
        "count":   res.ModifiedCount,
    }).Info()
    return r.GetByIDs(ctx, adminIDs...)
}

func (r *AdminsRepo) GetAll(ctx context.Context) ([]*repo.Admin, error) {
    cursor, err := r.this.Find(ctx, bson.D{}) // TODO: join roles
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "GetByIDs",
        }).Error(err)
        return nil, err
    }
    admins := make([]*repo.Admin, 0)
    if err := cursor.All(ctx, &admins); err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) GetByIDs(ctx context.Context, adminIDs ...int64) ([]*repo.Admin, error) {
    cursor, err := r.this.Find(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    }) // TODO: join roles
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "GetByIDs",
        }).Error(err)
        return nil, err
    }
    admins := make([]*repo.Admin, 0, len(adminIDs))
    if err := cursor.All(ctx, &admins); err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) GetByID(ctx context.Context, adminID int64) (*repo.Admin, error) {
    admins, err := r.GetByIDs(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if len(admins) != 1 {
        return nil, repo.AdminNotFound(adminID)
    }
    return admins[0], nil
}

func (r *AdminsRepo) HasScopesByID(ctx context.Context, adminID int64, scopes ...repo.Scope) (bool, error) {
    admin, err := r.GetByID(ctx, adminID)
    if err != nil {
        return false, err
    }
    if admin.Role == nil {
        return false, nil
    }
    for _, s := range scopes {
        if _, found := admin.Role.Scopes[s]; !found {
            return false, nil
        }
    }
    return true, nil
}

func (r *AdminsRepo) DeleteByIDs(ctx context.Context, adminIDs ...int64) error {
    res, err := r.this.DeleteMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "DeleteByIDs",
        }).Error(err)
        return err
    }
    log.WithFields(log.Fields{
        "context": "AdminsRepo",
        "action":  "DeleteByIDs",
        "count":   res.DeletedCount,
    }).Info()
    return nil
}
