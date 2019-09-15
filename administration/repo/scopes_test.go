package repo

import (
    "reflect"
    "strings"
    "testing"
)

func TestScope_Implies(t *testing.T) {
    tests := []struct {
        name string
        s    Scope
        want Scopes
    }{{
        name: string(StatsScope),
        s:    StatsScope,
        want: NewScopes(
            StatsScope,
        ),
    }, {
        name: string(AdminsReadScope),
        s:    AdminsReadScope,
        want: NewScopes(
            AdminsReadScope,
        ),
    }, {
        name: string(AdminsScope),
        s:    AdminsScope,
        want: NewScopes(
            AdminsScope,
            AdminsReadScope,
        ),
    },}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.s.Implies(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Implies() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestScopes_Has(t *testing.T) {
    type args struct {
        scopes []Scope
    }
    tests := []struct {
        name string
        s    Scopes
        args args
        want bool
    }{{
        name: scopesToString(
            StatsScope,
        ),
        s: NewScopes(
            StatsScope,
        ),
        args: args{scopes: []Scope{
            StatsScope,
        }},
        want: true,
    }, {
        name: scopesToString(
            StatsScope,
        ),
        s: NewScopes(
            StatsScope,
        ),
        args: args{scopes: []Scope{
            AdminsScope,
        }},
        want: false,
    }, {
        name: scopesToString(
            StatsScope,
        ),
        s: NewScopes(
            StatsScope,
        ),
        args: args{scopes: []Scope{
            AdminsReadScope,
        }},
        want: false,
    }, {
        name: scopesToString(
            AdminsReadScope,
        ),
        s: NewScopes(
            AdminsReadScope,
        ),
        args: args{scopes: []Scope{
            AdminsReadScope,
        }},
        want: true,
    }, {
        name: scopesToString(
            AdminsReadScope,
        ),
        s: NewScopes(
            AdminsReadScope,
        ),
        args: args{scopes: []Scope{
            AdminsScope,
        }},
        want: false,
    }, {
        name: scopesToString(
            AdminsReadScope,
        ),
        s: NewScopes(
            AdminsReadScope,
        ),
        args: args{scopes: []Scope{
            StatsScope,
        }},
        want: false,
    }, {
        name: scopesToString(
            AdminsScope,
        ),
        s: NewScopes(
            AdminsScope,
        ),
        args: args{scopes: []Scope{
            AdminsScope,
            AdminsReadScope,
        }},
        want: true,
    }, {
        name: scopesToString(
            StatsScope,
            AdminsScope,
        ),
        s: NewScopes(
            StatsScope,
            AdminsScope,
        ),
        args: args{scopes: []Scope{
            StatsScope,
            AdminsScope,
            AdminsReadScope,
        }},
        want: true,
    },}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.s.Has(tt.args.scopes...); got != tt.want {
                t.Errorf("Has() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestScopes_Implies(t *testing.T) {
    tests := []struct {
        name string
        s    Scopes
        want Scopes
    }{{
        name: scopesToString(
            StatsScope,
        ),
        s: NewScopes(
            StatsScope,
        ),
        want: NewScopes(
            StatsScope,
        ),
    }, {
        name: scopesToString(
            AdminsReadScope,
        ),
        s: NewScopes(
            AdminsReadScope,
        ),
        want: NewScopes(
            AdminsReadScope,
        ),
    }, {
        name: scopesToString(
            AdminsScope,
        ),
        s: NewScopes(
            AdminsScope,
        ),
        want: NewScopes(
            AdminsScope,
            AdminsReadScope,
        ),
    }, {
        name: scopesToString(
            AdminsScope,
            AdminsReadScope,
        ),
        s: NewScopes(
            AdminsScope,
            AdminsReadScope,
        ),
        want: NewScopes(
            AdminsScope,
            AdminsReadScope,
        ),
    }, {
        name: scopesToString(
            AdminsScope,
            StatsScope,
        ),
        s: NewScopes(
            AdminsScope,
            StatsScope,
        ),
        want: NewScopes(
            AdminsScope,
            AdminsReadScope,
            StatsScope,
        ),
    }, {
        name: scopesToString(
            AdminsReadScope,
            StatsScope,
        ),
        s: NewScopes(
            AdminsReadScope,
            StatsScope,
        ),
        want: NewScopes(
            AdminsReadScope,
            StatsScope,
        ),
    }, {
        name: scopesToString(
            AdminsReadScope,
            AdminsScope,
            StatsScope,
        ),
        s: NewScopes(
            AdminsReadScope,
            AdminsScope,
            StatsScope,
        ),
        want: NewScopes(
            AdminsReadScope,
            AdminsScope,
            StatsScope,
        ),
    },}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.s.Implies(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Implies() = %v, want %v", got, tt.want)
            }
        })
    }
}

func scopesToString(scopes ...Scope) string {
    s := make([]string, 0, len(scopes))
    for _, scope := range scopes {
        s = append(s, string(scope))
    }
    return strings.Join(s, ",")
}
