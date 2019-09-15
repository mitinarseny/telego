package repo

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
