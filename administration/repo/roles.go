package repo

import "context"

type Scope string

const (
    AdminsScope Scope = "admins"
    StatsScope  Scope = "stats"

    SuperUserRoleName = "superuser"
)

var (
    SuperuserRole = NewRole(SuperUserRoleName,
        AdminsScope,
        StatsScope,
    )
)

type Role struct {
    Name   string             `bson:"_id,omitempty"`
    Scopes map[Scope]struct{} `bson:"scopes,omitempty"`
}

func NewRole(name string, scopes ...Scope) *Role{
    r := &Role{
        Name: name,
        Scopes: make(map[Scope]struct{}, len(scopes)),
    }
    for _, s := range scopes {
        r.Scopes[s] = struct{}{}
    }
    return r
}

type RolesRepo interface {
    Create(ctx context.Context, roles ...*Role) ([]*Role, error)
    CreateIfNotExists(ctx context.Context, roles ...*Role) ([]*Role, error)
    GetByNames(ctx context.Context, names ...string) ([]*Role, error)
    AddScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    SetScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    DeleteByNames(ctx context.Context, names ...string) error
}
