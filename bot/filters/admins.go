package filters

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

type isAdmin struct {
    msgParent      MsgFilter
    callbackParent CallbackFilter

    admins repo.AdminsRepo
}

func (f *isAdmin) isAdmin(adminID int64) (bool, error) {
    _, err := f.admins.GetByID(context.Background(), adminID)
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

func (f *isAdmin) FilterMsg(m *tb.Message) (bool, error) {
    if f.msgParent != nil {
        if passed, err := f.msgParent.FilterMsg(m); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return f.isAdmin(int64(m.Sender.ID))
}

func (f *isAdmin) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        if passed, err := f.callbackParent.FilterCallback(c); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return f.isAdmin(int64(c.Sender.ID))
}

func (f *isAdmin) HasScopes(scopes ...repo.Scope) *hasScopes {
    return &hasScopes{
        msgParent:      f,
        callbackParent: f,
        scopes:         scopes,
    }
}

type hasScopes struct {
    callbackParent CallbackFilter
    msgParent      MsgFilter

    admins repo.AdminsRepo
    scopes []repo.Scope
}

func (f *hasScopes) hasScopes(adminID int64) (bool, error) {
    has, err := f.admins.HasScopesByID(context.Background(), adminID, f.scopes...)
    if err != nil {
        return false, err
    }
    return has, nil
}

func (f *hasScopes) FilterMsg(m *tb.Message) (bool, error) {
    if passed, err := f.msgParent.FilterMsg(m); err != nil {
        return false, err
    } else if !passed {
        return false, nil
    }
    return f.hasScopes(int64(m.Sender.ID))
}

func (f *hasScopes) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        if passed, err := f.callbackParent.FilterCallback(c); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return f.hasScopes(int64(c.Sender.ID))
}

type isAdminWithScopes struct {
    callbackParent CallbackFilter
    msgParent      MsgFilter

    admins repo.AdminsRepo
    scopes []repo.Scope
}

func (f *isAdminWithScopes) isAdminWithScopes(adminID int64) (bool, error) {
    has, err := f.admins.HasScopesByID(context.Background(), adminID, f.scopes...)
    if err != nil {
        switch err.(type) {
        case repo.AdminNotFound:
            return false, nil
        default:
            return false, err
        }
    }
    return has, nil
}

func (f *isAdminWithScopes) FilterMsg(m *tb.Message) (bool, error) {
    if passed, err := f.msgParent.FilterMsg(m); err != nil {
        return false, err
    } else if !passed {
        return false, nil
    }
    return f.isAdminWithScopes(int64(m.Sender.ID))
}

func (f *isAdminWithScopes) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        if passed, err := f.callbackParent.FilterCallback(c); err != nil {
            return false, err
        } else if !passed {
            return false, nil
        }
    }
    return f.isAdminWithScopes(int64(c.Sender.ID))
}
