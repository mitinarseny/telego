package admins

import (
    "context"

    "go.mongodb.org/mongo-driver/bson"
)

const SuperUserRoleName = "superuser"

var SuperuserRole = NewRole(SuperUserRoleName,
    AdminsScope,
    StatsScope,
)

type Role struct {
    Name   string `bson:"name,omitempty"`
    Scopes Scopes `bson:"scopes,omitempty"`
}

type roleBson struct {
    Name   string  `bson:"name,omitempty"`
    Scopes []Scope `bson:"scopes,omitempty"`
}

func (r *Role) MarshalBSON() ([]byte, error) {
    aux := roleBson{
        Name:   r.Name,
        Scopes: make([]Scope, 0, len(r.Scopes)),
    }
    for sc := range r.Scopes {
        aux.Scopes = append(aux.Scopes, sc)
    }
    return bson.Marshal(&aux)
}

func (r *Role) UnmarshalBSON(data []byte) error {
    var aux roleBson
    if err := bson.Unmarshal(data, &aux); err != nil {
        return err
    }
    r.Name = aux.Name
    r.Scopes = make(Scopes, len(aux.Scopes))
    for _, sc := range aux.Scopes {
        r.Scopes[sc] = struct{}{}
    }
    return nil
}

func (r *Role) HasScopes(scopes ...Scope) bool {
    return r.Scopes.Has(scopes...)
}

func NewRole(name string, scopes ...Scope) *Role {
    r := &Role{
        Name:   name,
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
    GetAll(ctx context.Context) ([]*Role, error)
    GetByName(ctx context.Context, name string) (*Role, error)
    GetByNames(ctx context.Context, names ...string) ([]*Role, error)
    AddScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    SetScopes(ctx context.Context, scopes []Scope, names ...string) ([]*Role, error)
    DeleteByNames(ctx context.Context, names ...string) error
}
