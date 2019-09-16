package filters

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

type AdminsOnly struct {
    Admins repo.AdminsRepo
}

func (f *AdminsOnly) Filter(m *tb.Message) (bool, error) {
    _, err := f.Admins.GetByID(context.Background(), int64(m.Sender.ID))
    if err != nil {
        switch err.(type) {
        case repo.AdminNotFound:
            return false, nil
        default:
            return false, err
        }
    }
    return true, nil
}

func (f *AdminsOnly) WithScopes(scopes ...repo.Scope) *onlyAdminsWithScopes {
    return &onlyAdminsWithScopes{
        AdminsOnly: f,
        scopes:     scopes,
    }
}

type onlyAdminsWithScopes struct {
    *AdminsOnly
    scopes []repo.Scope
}

func (f *onlyAdminsWithScopes) Filter(m *tb.Message) (bool, error) {
    has, err := f.Admins.HasScopesByID(context.Background(), int64(m.Sender.ID), f.scopes...)
    if err != nil {
        return false, err
    }
    return has, nil
}
