package repo

import "context"

type Scope string

const (
    SuperUserRoleName = "superuser"
)

type Role struct {
    Name   string             `bson:"_id,omitempty"`
    Scopes map[Scope]struct{} `bson:"scopes,omitempty"`
}

type RolesRepo interface {
    Create(ctx context.Context, roles ...*Role) ([]*Role, error)
    CreateIfNotExists(ctx context.Context, roles ...*Role) ([]*Role, error)
    GetByNames(ctx context.Context, names ...string) ([]*Role, error)
    AddScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    SetScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    DeleteByNames(ctx context.Context, names ...string) error
}
