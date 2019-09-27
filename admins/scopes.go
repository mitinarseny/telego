package admins

import (
    "go.mongodb.org/mongo-driver/bson"
)

const (
    AdminsScope     Scope = "admins"
    AdminsReadScope Scope = "admins.read"
    StatsScope      Scope = "stats"
)

var scopeHierarchy = map[Scope][]Scope{
    AdminsScope: {
        AdminsReadScope,
    },
}

type Scope string

func (s Scope) Implies() Scopes {
    scopes := Scopes{s: {}}
    for _, scopeInHier := range scopeHierarchy[s] {
        for sc := range scopeInHier.Implies() {
            scopes[sc] = struct{}{}
        }
    }
    return scopes
}

type Scopes map[Scope]struct{}

func (s Scopes) MarshalBSON() ([]byte, error) {
    aux := make(bson.A, 0, len(s))
    for sc := range s {
        aux = append(aux, sc)
    }
    return bson.Marshal(aux)
}

func (s Scopes) UnmarshalBSON(data []byte) error {
    aux := make([]Scope, 0)
    if err := bson.Unmarshal(data, &aux); err != nil {
        return err
    }
    for _, sc := range aux {
        s[sc] = struct{}{}
    }
    return nil
}

func NewScopes(scopes ...Scope) Scopes {
    res := make(Scopes, len(scopes))
    for _, s := range scopes {
        res[s] = struct{}{}
    }
    return res
}

func (s Scopes) Implies() Scopes {
    scopes := make(Scopes)
    for scope := range s {
        for implSc := range scope.Implies() {
            scopes[implSc] = struct{}{}
        }
    }
    return scopes
}

func (s Scopes) Has(scopes ...Scope) bool {
    implied := s.Implies()
    for _, sc := range scopes {
        if _, found := implied[sc]; !found {
            return false
        }
    }
    return true
}
