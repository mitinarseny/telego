package filters

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (f *Filters) OnlyFromSuperUsers(h MessageHandler) MessageHandler {
    return func(m *tb.Message) error {
        if sender := m.Sender; sender != nil {
            admin, err := f.storage.Admins.GetByID(context.Background(), int64(sender.ID))
            if err != nil {
                switch err.(type) {
                case repo.AdminNotFound:
                    return nil
                default:
                    return err
                }
            }
            if admin.Role == nil {
                return nil
            }
            if admin.Role.Name != repo.SuperUserRoleName {
                return nil
            }
            return h(m)
        }
        return nil
    }
}
