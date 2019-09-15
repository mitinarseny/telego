package bot

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) superusersOnly(m *tb.Message) bool {
    if m.Sender == nil {
        return false
    }
    admin, err := b.storage.Admins.GetByID(context.Background(), int64(m.Sender.ID))
    if err != nil {
        log.WithFields(log.Fields{
            "context": "BOT.superusersOnly",
            "action": "GetByID",
            "id": m.Sender.ID,
        }).Error(err)
        return false
    }
    if admin.Role == nil {
        log.WithFields(log.Fields{
            "context": "BOT.superusersOnly",
            "id": admin.ID,
        }).Error("admin do not has a role")
        return false
    }
    if admin.Role.Name != repo.SuperUserRoleName {
        return false
    }
    return true
}
