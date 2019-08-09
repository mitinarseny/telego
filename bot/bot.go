package bot

import (
    "github.com/mitinarseny/telego/bot/handlers"
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

func Configure(b *tb.Bot) (*tb.Bot, error) {
    h := handlers.Handler{Bot: b}
    b.Handle("/hello", withLogMsg(h.HandleHello))
    return b, nil
}

type MessageHandler func(*tb.Message) error

func withLogMsg(h MessageHandler) func(*tb.Message) {
    return func(m *tb.Message) {
        if err := h(m); err != nil {
            log.WithFields(log.Fields{
                "context": "BOT",
                "handler":       h,
            }).Error(err)
        }
    }
}
