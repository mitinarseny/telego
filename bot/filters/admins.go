package filters

import (
    "context"

    "github.com/mitinarseny/telego/admins"
    tb "gopkg.in/tucnak/telebot.v2"
)

type isAdmin struct {
    msgParent      MsgFilter
    callbackParent CallbackFilter

    admins admins.AdminsRepo
}

func (f *isAdmin) isAdmin(adminID int64) (bool, error) {
    if _, err := f.admins.GetByID(context.Background(), adminID); err != nil {
        switch err.(type) {
        case admins.AdminNotFound:
            return false, nil
        default:
            return false, err
        }
    }
    return true, nil
}

func (f *isAdmin) FilterMsg(m *tb.Message) (bool, error) {
    if f.msgParent != nil {
        switch passed, err := f.msgParent.FilterMsg(m); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.isAdmin(int64(m.Sender.ID))
}

func (f *isAdmin) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        switch passed, err := f.callbackParent.FilterCallback(c) ;{
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.isAdmin(int64(c.Sender.ID))
}

func (f *isAdmin) HasScopes(scopes ...admins.Scope) *hasScopes {
    return &hasScopes{
        msgParent:      f,
        callbackParent: f,
        scopes:         scopes,
    }
}

type hasScopes struct {
    callbackParent CallbackFilter
    msgParent      MsgFilter

    admins admins.AdminsRepo
    scopes []admins.Scope
}

func (f *hasScopes) hasScopes(adminID int64) (bool, error) {
    has, err := f.admins.HasScopesByID(context.Background(), adminID, f.scopes...)
    if err != nil {
        return false, err
    }
    return has, nil
}

func (f *hasScopes) FilterMsg(m *tb.Message) (bool, error) {
    if f.msgParent != nil {
        switch passed, err := f.msgParent.FilterMsg(m); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.hasScopes(int64(m.Sender.ID))
}

func (f *hasScopes) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        switch passed, err := f.callbackParent.FilterCallback(c); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.hasScopes(int64(c.Sender.ID))
}

type isAdminWithScopes struct {
    callbackParent CallbackFilter
    msgParent      MsgFilter

    admins admins.AdminsRepo
    scopes []admins.Scope
}

func (f *isAdminWithScopes) isAdminWithScopes(adminID int64) (bool, error) {
    has, err := f.admins.HasScopesByID(context.Background(), adminID, f.scopes...)
    if err != nil {
        switch err.(type) {
        case admins.AdminNotFound:
            return false, nil
        default:
            return false, err
        }
    }
    return has, nil
}

func (f *isAdminWithScopes) FilterMsg(m *tb.Message) (bool, error) {
    if f.msgParent != nil {
        switch passed, err := f.msgParent.FilterMsg(m); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.isAdminWithScopes(int64(m.Sender.ID))
}

func (f *isAdminWithScopes) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        switch passed, err := f.callbackParent.FilterCallback(c); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.isAdminWithScopes(int64(c.Sender.ID))
}
