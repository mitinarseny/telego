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

func withLogMsg(handler func(*tb.Message) error) func(message *tb.Message) {
    return func(m *tb.Message) {
        if err := handler(m); err != nil {
            log.WithFields(log.Fields{
                "context": "BOT",
                "handler": handler,
            }).Error(err)
        }
    }
}
