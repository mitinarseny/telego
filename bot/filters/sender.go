package filters

import (
    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

func WithSender() *sender {
    return &sender{}
}

type sender struct {
    msgParent      MsgFilter
    callbackParent CallbackFilter
}

func (f *sender) FilterMsg(m *tb.Message) (bool, error) {
    if f.msgParent != nil {
        if passed, err := f.msgParent.FilterMsg(m); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return m.Sender != nil, nil
}

func (f *sender) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        if passed, err := f.callbackParent.FilterCallback(c); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return c.Sender != nil, nil
}

func (f *sender) IsAdmin(r repo.AdminsRepo) *isAdmin {
    return &isAdmin{
        msgParent:      f,
        callbackParent: f,
        admins:         r,
    }
}

func (f *sender) IsAdminWithScopes(r repo.AdminsRepo, scopes ...repo.Scope) *isAdminWithScopes {
    return &isAdminWithScopes{
        callbackParent:f,
        msgParent:f,

        admins:r,
        scopes: scopes,
    }
}
