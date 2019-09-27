package filters

import (
    "github.com/mitinarseny/telego/admins"
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
        switch passed, err := f.msgParent.FilterMsg(m); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return m.Sender != nil, nil
}

func (f *sender) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        switch passed, err := f.callbackParent.FilterCallback(c); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return c.Sender != nil, nil
}

func (f *sender) IsAdmin(r admins.AdminsRepo) *isAdmin {
    return &isAdmin{
        msgParent:      f,
        callbackParent: f,
        admins:         r,
    }
}

func (f *sender) IsAdminWithScopes(r admins.AdminsRepo, scopes ...admins.Scope) *isAdminWithScopes {
    return &isAdminWithScopes{
        callbackParent: f,
        msgParent:      f,

        admins: r,
        scopes: scopes,
    }
}
