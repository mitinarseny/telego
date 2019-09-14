package handlers

import (
    "context"
    "fmt"

    "github.com/mitinarseny/telego/administration/repo"
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (h *Handler) HandleAdmins(m *tb.Message) error {
    admins, err := h.Storage.Admins.GetAll(context.Background())
    if err != nil {
        return err
    }
    inlineKeys := h.makeAdminsInlineKeyboard(admins)
    _, err = h.Bot.Send(m.Sender, "Here are admins", &tb.ReplyMarkup{
        InlineKeyboard: inlineKeys,
    })
    return err
}

func (h *Handler) makeAdminsInlineKeyboard(admins []*repo.Admin) [][]tb.InlineButton {
    keyboard := make([][]tb.InlineButton, 0, len(admins)/2+len(admins)%2)
    for i, admin := range admins {
        if i%2 == 0 {
            keyboard = append(keyboard, make([]tb.InlineButton, 0, 2))
        }
        btn := tb.InlineButton{
            Unique: fmt.Sprintf("%d", admin.ID), // TODO: more unique
            Text:   fmt.Sprintf("%d (%s)", admin.ID, admin.Role.Name),
        }
        h.Bot.Handle(&btn, func(c *tb.Callback) {
            log.WithFields(log.Fields{
                "context":       "BOT",
                "handler":       "HandleAdminsCallback",
                "callback_data": c.Data,
            }).Info()
            if err := h.Bot.Respond(c, &tb.CallbackResponse{}); err != nil {
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
