package repo

import (
    "context"
    "fmt"
    "strconv"
)

type AdminNotFound int64

func (e AdminNotFound) Error() string {
    return fmt.Sprintf("Admin %q not found", strconv.FormatInt(int64(e), 10))
}

type Admin struct {
    ID   int64 `bson:"_id, omitempty"`
    Role *Role `bson:"role,omitempty"`
}

type AdminsRepo interface {
    Create(ctx context.Context, admins ...*Admin) ([]*Admin, error)
    CreateIfNotExists(ctx context.Context, admins ...*Admin) ([]*Admin, error)
    ChangeRoleByIDs(ctx context.Context, role *Role, adminIDs ...int64) ([]*Admin, error)
    GetAll(ctx context.Context) ([]*Admin, error)
    GetByID(ctx context.Context, adminID int64) (*Admin, error)
    GetByIDs(ctx context.Context, adminIDs ...int64) ([]*Admin, error)
    HasScopesByID(ctx context.Context, adminID int64, scopes ...Scope) (bool, error)
    DeleteByIDs(ctx context.Context, adminIDs ...int64) error
}
