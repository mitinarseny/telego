package mongo

import (
    "context"
    "errors"
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/bot/repo/administration"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    adminsCollectionName = "admins"
)

type AdminsRepo struct {
    this  *mongo.Collection
    roles administration.RolesRepo
}

func NewAdminsRepo(db *mongo.Database, roles administration.RolesRepo) *AdminsRepo {
    return &AdminsRepo{
        this:  db.Collection(adminsCollectionName),
        roles: roles,
    }
}

func (r *AdminsRepo) Create(ctx context.Context, admins ...*administration.Admin) ([]*administration.Admin, error) {
    models := make([]interface{}, 0, len(admins))
    roles := make([]*administration.Role, 0, len(admins))
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
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return admins, nil
}

type adminModel struct {
    ID            int64                                   `bson:"_id, omitempty"`
    RoleName      string                                  `bson:"roleName,omitempty"`
    Notifications administration.NotificationsPreferences `bson:"notifications,omitempty"`
}

func (r *AdminsRepo) CreateIfNotExists(ctx context.Context, admins ...*administration.Admin) ([]*administration.Admin, error) {
    models := make([]mongo.WriteModel, 0, len(admins))
    roles := make([]*administration.Role, 0, len(admins))
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
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) AssignRoleByID(ctx context.Context, roleName string, adminID int64) (*administration.Admin, error) {
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

func (r *AdminsRepo) AssignRoleByIDs(ctx context.Context, roleName string, adminIDs ...int64) ([]*administration.Admin, error) {
    role, err := r.roles.GetByName(ctx, roleName)
    if err != nil {
        return nil, err
    }
    _, err = r.this.UpdateMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    }, bson.D{
        {"$set", bson.D{
            {"roleName", role.Name},
        }},
    })
    if err != nil {
        return nil, err
    }
    return r.GetByIDs(ctx, adminIDs...)
}

func (r *AdminsRepo) GetAll(ctx context.Context) ([]*administration.Admin, error) {
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
        return nil, err
    }
    admins := make([]*administration.Admin, 0)
    if err := cursor.All(ctx, &admins); err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) GetAllShouldBeNotifiedAbout(ctx context.Context,
    notificationType administration.NotificationType) ([]*administration.Admin, error) {
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{
            {"$match", bson.D{
                {"notifications", bson.D{
                    {string(notificationType), bson.D{
                        {"$type", "array"},
                        {"$not", bson.D{
                            {"$size", 0},
                        }},
                    }},
                }},
            }},
        },
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
        return nil, err
    }
    admins := make([]*administration.Admin, 0)
    if err := cursor.All(ctx, &admins); err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) GetByIDs(ctx context.Context, adminIDs ...int64) ([]*administration.Admin, error) {
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
        return nil, err
    }
    admins := make([]*administration.Admin, 0, len(adminIDs))
    if err := cursor.All(ctx, &admins); err != nil {
        return nil, err
    }
    return admins, nil
}

func (r *AdminsRepo) GetByID(ctx context.Context, adminID int64) (*administration.Admin, error) {
    admins, err := r.GetByIDs(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if len(admins) != 1 {
        return nil, administration.AdminNotFound(adminID)
    }
    return admins[0], nil
}

func (r *AdminsRepo) GetRoleByID(ctx context.Context, adminID int64) (*administration.Role, error) {
    admin, err := r.GetByID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if admin.Role == nil {
        return nil, errors.New("empty role")
    }
    return admin.Role, nil
}

func (r *AdminsRepo) HasScopesByID(ctx context.Context, adminID int64, scopes ...administration.Scope) (bool, error) {
    role, err := r.GetRoleByID(ctx, adminID)
    if err != nil {
        return false, err
    }
    return role.HasScopes(scopes...), nil
}

func (r *AdminsRepo) DeleteByIDs(ctx context.Context, adminIDs ...int64) error {
    _, err := r.this.DeleteMany(ctx, bson.D{
        {"_id", bson.D{
            {"$in", adminIDs},
        }},
    })
    return err
}
