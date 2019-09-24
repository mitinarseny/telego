package mongo

import (
    "context"
    "errors"
    "fmt"
    "strconv"

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
    models := make([]interface{}, 0, len(admins))
    roles := make([]*repo.Role, 0, len(admins))
    for _, admin := range admins {
        models = append(models, &adminModel{
            ID:       admin.ID,
            RoleName: admin.Role.Name,
        })
        roles = append(roles, admin.Role)
    }
    if _, err := r.roles.CreateIfNotExists(ctx, roles...); err != nil {
        return nil, err
    }
    res, err := r.this.InsertMany(ctx, models)
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

type adminModel struct {
    ID       int64  `bson:"_id, omitempty"`
    RoleName string `bson:"roleName,omitempty"`
}

func (r *AdminsRepo) CreateIfNotExists(ctx context.Context, admins ...*repo.Admin) ([]*repo.Admin, error) {
    models := make([]mongo.WriteModel, 0, len(admins))
    roles := make([]*repo.Role, 0, len(admins))
    for _, admin := range admins {
        roles = append(roles, admin.Role)
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"_id", admin.ID},
            }).SetUpdate(&adminModel{
                ID:       admin.ID,
                RoleName: admin.Role.Name,
            }).SetUpsert(true))
    }
    if _, err := r.roles.CreateIfNotExists(ctx, roles...); err != nil {
        return nil, err
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

func (r *AdminsRepo) AssignRoleByID(ctx context.Context, roleName string, adminID int64) (*repo.Admin, error) {
    admins, err := r.AssignRoleByIDs(ctx, roleName, adminID)
    if err != nil {
        return nil, err
    }
    switch {
    case len(admins) == 0:
        return nil, errors.New(fmt.Sprintf("admin %q not found", strconv.FormatInt(adminID, 10)))
    case len(admins) > 1:
        return nil, errors.New(fmt.Sprintf("more than one admin with ID %q found", strconv.FormatInt(adminID, 10)))
    }
    return admins[0], nil
}

func (r *AdminsRepo) AssignRoleByIDs(ctx context.Context, roleName string, adminIDs ...int64) ([]*repo.Admin, error) {
    role, err := r.roles.GetByName(ctx, roleName)
    if err != nil {
        return nil, err
    }
    res, err := r.this.UpdateMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    }, bson.D{
        {"$set", bson.D{
            {"roleName", role.Name},
        }},
    })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "action":  "AssignRoleByIDs",
        }).Error(err)
        return nil, err
    }
    log.WithFields(log.Fields{
        "context": "UpdatesRepo",
        "action":  "AssignRoleByIDs",
        "count":   res.ModifiedCount,
    }).Info()
    return r.GetByIDs(ctx, adminIDs...)
}

func (r *AdminsRepo) GetAll(ctx context.Context) ([]*repo.Admin, error) {
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{
            {"$lookup", bson.D{
                {"from", rolesCollectionName},
                {"localField", "roleName"},
                {"foreignField", "_id"},
                {"as", "role"},
            }},
        }, bson.D{
            {"$project", bson.D{
                {"role", bson.D{
                    {"$arrayElemAt", bson.A{"$role", 0}},
                }},
            }},
        },
    })
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
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{{"$match", bson.D{
            {"_id", bson.D{
                {"$in", adminIDs},
            }},
        }}}, bson.D{
            {"$lookup", bson.D{
                {"from", rolesCollectionName},
                {"localField", "roleName"},
                {"foreignField", "_id"},
                {"as", "role"},
            }},
        }, bson.D{
            {"$project", bson.D{
                {"role", bson.D{
                    {"$arrayElemAt", bson.A{"$role", 0}},
                }},
            }},
        },
    })
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

func (r *AdminsRepo) GetRoleByID(ctx context.Context, adminID int64) (*repo.Role, error) {
    admin, err := r.GetByID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if admin.Role == nil {
        log.WithFields(log.Fields{
            "context": "AdminsRepo",
            "id":      admin.ID,
        }).Error("role is empty")
        return nil, errors.New("empty role")
    }
    return admin.Role, nil
}

func (r *AdminsRepo) HasScopesByID(ctx context.Context, adminID int64, scopes ...repo.Scope) (bool, error) {
    role, err := r.GetRoleByID(ctx, adminID)
    if err != nil {
        return false, err
    }
    return role.HasScopes(scopes...), nil
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
