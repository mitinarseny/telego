package bot

import (
    "context"

    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onlyAdminsWithScopes(scopes ...repo.Scope) MsgFilter {
    return func(m *tb.Message) (bool, error) {
        if m.Sender == nil {
            return false, nil
        }
        has, err := b.storage.Admins.HasScopesByID(context.Background(), int64(m.Sender.ID), scopes...)
        if err != nil {
            return false, err
        }
        return has, nil
    }
}
