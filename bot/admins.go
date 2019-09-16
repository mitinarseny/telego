package bot

import (
    "context"
    "fmt"

    "github.com/mitinarseny/telego/administration/repo"
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) handleAdmins(m *tb.Message) error {
    admins, err := b.storage.Admins.GetAll(context.Background())
    if err != nil {
        return err
    }
    inlineKeys := b.makeAdminsInlineKeyboard(admins)
    _, err = b.tg.Send(m.Sender, "Here are admins", &tb.ReplyMarkup{
        InlineKeyboard: inlineKeys,
    })
    return err
}

func (b *Bot) makeAdminsInlineKeyboard(admins []*repo.Admin) [][]tb.InlineButton {
    keyboard := make([][]tb.InlineButton, 0, len(admins)/2+len(admins)%2)
    for i, admin := range admins {
        if i%2 == 0 {
            keyboard = append(keyboard, make([]tb.InlineButton, 0, 2))
        }
        btn := tb.InlineButton{
            Unique: fmt.Sprintf("%d", admin.ID), // TODO: more unique
            Text:   fmt.Sprintf("%d (%s)", admin.ID, admin.Role.Name),
        }
        b.tg.Handle(&btn, func(c *tb.Callback) {
            log.WithFields(log.Fields{
                "context":       "BOT",
                "handler":       "HandleAdminsCallback",
                "callback_data": c.Data,
            }).Info()
            if err := b.tg.Respond(c, &tb.CallbackResponse{}); err != nil {
                log.WithFields(log.Fields{
                    "context": "BOT",
                    "handler": "HandleAdminsCallback",
                }).Error(err)
            }
        })
        keyboard[i/2] = append(keyboard[i/2], btn)
    }
    return keyboard
}