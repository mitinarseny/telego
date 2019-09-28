package mongo

import (
    "context"
    "fmt"
    "strconv"

    "github.com/pkg/errors"

    "github.com/mitinarseny/telego/admins"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    adminsCollectionName = "admins"
)

type AdminsRepo struct {
    this  *mongo.Collection
    roles admins.RolesRepo
}

func NewAdminsRepo(db *mongo.Database, roles admins.RolesRepo) (*AdminsRepo, error) {
    r := &AdminsRepo{
        this:  db.Collection(adminsCollectionName),
        roles: roles,
    }
    if _, err := r.this.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{"id", 1}},
        Options: options.Index().SetUnique(true),
    }); err != nil {
        return nil, errors.Wrapf(err, "unable to create index on %s", adminsCollectionName)
    }
    return r, nil
}

func (r *AdminsRepo) Create(ctx context.Context, adms ...*admins.Admin) ([]*admins.Admin, error) {
    models := make([]interface{}, 0, len(adms))
    roles := make([]*admins.Role, 0, len(adms))
    for _, admin := range adms {
        models = append(models, adminModelFromRepo(admin))
        roles = append(roles, admin.Role)
    }
    if _, err := r.roles.CreateIfNotExists(ctx, roles...); err != nil {
        return nil, err
    }
    _, err := r.this.InsertMany(ctx, models)
    if err != nil {
        return nil, err
    }
    return adms, nil
}

type adminModel struct {
    ID            int64                 `bson:"id,omitempty"`
    RoleName      string                `bson:"roleName,omitempty"`
    Notifications *admins.Notifications `bson:"notifications,omitempty"`
}

func adminModelFromRepo(a *admins.Admin) *adminModel {
    return &adminModel{
        ID:            a.ID,
        RoleName:      a.Role.Name,
        Notifications: a.Notifications,
    }
}

func (r *AdminsRepo) CreateIfNotExists(ctx context.Context, adms ...*admins.Admin) ([]*admins.Admin, error) {
    models := make([]mongo.WriteModel, 0, len(adms))
    roles := make([]*admins.Role, 0, len(adms))
    for _, admin := range adms {
        roles = append(roles, admin.Role)
        models = append(models,
            mongo.NewUpdateOneModel().SetFilter(bson.D{
                {"id", admin.ID},
            }).SetUpdate(adminModelFromRepo(admin)).SetUpsert(true))
    }
    if _, err := r.roles.CreateIfNotExists(ctx, roles...); err != nil {
        return nil, err
    }
    _, err := r.this.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
    if err != nil {
        return nil, err
    }
    return adms, nil
}

func (r *AdminsRepo) AssignRoleByID(ctx context.Context, roleName string, adminID int64) (*admins.Admin, error) {
    adms, err := r.AssignRoleByIDs(ctx, roleName, adminID)
    if err != nil {
        return nil, err
    }
    switch {
    case len(adms) == 0:
        return nil, errors.New(fmt.Sprintf("admin %q not found", strconv.FormatInt(adminID, 10)))
    case len(adms) > 1:
        return nil, errors.New(fmt.Sprintf("more than one admin with ID %q found", strconv.FormatInt(adminID, 10)))
    }
    return adms[0], nil
}

func (r *AdminsRepo) AssignRoleByIDs(ctx context.Context, roleName string, adminIDs ...int64) ([]*admins.Admin, error) {
    role, err := r.roles.GetByName(ctx, roleName)
    if err != nil {
        return nil, err
    }
    _, err = r.this.UpdateMany(ctx, bson.D{
        {"id", bson.D{
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

func (r *AdminsRepo) GetAll(ctx context.Context) ([]*admins.Admin, error) {
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{
            {"$lookup", bson.D{
                {"from", rolesCollectionName},
                {"localField", "roleName"},
                {"foreignField", "name"},
                {"as", "role"},
            }},
        }, bson.D{
            {"$project", bson.D{
                {"role", bson.D{
                    {"$arrayElemAt", bson.A{"$role", 0}},
                }},
                {"notifications", 1},
                {"id", 1},
            }},
        },
    })
    if err != nil {
        return nil, err
    }
    adms := make([]*admins.Admin, 0)
    if err := cursor.All(ctx, &adms); err != nil {
        return nil, err
    }
    return adms, nil
}

func (r *AdminsRepo) GetAllShouldBeNotifiedAbout(ctx context.Context,
    notificationType admins.NotificationType) ([]*admins.Admin, error) {
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{
            {"$match", bson.D{
                {"notifications." + string(notificationType), bson.D{
                    {"$type", "array"},
                    {"$not", bson.D{
                        {"$size", 0},
                    }},
                }},
            }},
        },
        bson.D{
            {"$lookup", bson.D{
                {"from", rolesCollectionName},
                {"localField", "roleName"},
                {"foreignField", "name"},
                {"as", "role"},
            }},
        }, bson.D{
            {"$project", bson.D{
                {"role", bson.D{
                    {"$arrayElemAt", bson.A{"$role", 0}},
                }},
                {"notifications", 1},
                {"id", 1},
            }},
        },
    })
    if err != nil {
        return nil, err
    }
    adms := make([]*admins.Admin, 0)
    if err := cursor.All(ctx, &adms); err != nil {
        return nil, err
    }
    return adms, nil
}

func (r *AdminsRepo) GetByIDs(ctx context.Context, adminIDs ...int64) ([]*admins.Admin, error) {
    cursor, err := r.this.Aggregate(ctx, bson.A{
        bson.D{{"$match", bson.D{
            {"id", bson.D{
                {"$in", adminIDs},
            }},
        }}}, bson.D{
            {"$lookup", bson.D{
                {"from", rolesCollectionName},
                {"localField", "roleName"},
                {"foreignField", "name"},
                {"as", "role"},
            }},
        }, bson.D{
            {"$project", bson.D{
                {"role", bson.D{
                    {"$arrayElemAt", bson.A{"$role", 0}},
                }},
                {"notifications", 1},
                {"id", 1},
            }},
        },
    })
    if err != nil {
        return nil, err
    }
    adms := make([]*admins.Admin, 0, len(adminIDs))
    if err := cursor.All(ctx, &adms); err != nil {
        return nil, err
    }
    return adms, nil
}

func (r *AdminsRepo) GetByID(ctx context.Context, adminID int64) (*admins.Admin, error) {
    adms, err := r.GetByIDs(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if len(adms) != 1 {
        return nil, admins.AdminNotFound(adminID)
    }
    return adms[0], nil
}

func (r *AdminsRepo) GetRoleByID(ctx context.Context, adminID int64) (*admins.Role, error) {
    admin, err := r.GetByID(ctx, adminID)
    if err != nil {
        return nil, err
    }
    if admin.Role == nil {
        return nil, errors.New("empty role")
    }
    return admin.Role, nil
}

func (r *AdminsRepo) HasScopesByID(ctx context.Context, adminID int64, scopes ...admins.Scope) (bool, error) {
    role, err := r.GetRoleByID(ctx, adminID)
    if err != nil {
        return false, err
    }
    return role.HasScopes(scopes...), nil
}

func (r *AdminsRepo) DeleteByIDs(ctx context.Context, adminIDs ...int64) error {
    _, err := r.this.DeleteMany(ctx, bson.D{
        {"id", bson.D{
            {"$in", adminIDs},
        }},
    })
    return err
}
