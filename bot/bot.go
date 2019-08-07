package bot

import (
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

type UpdatesLogger interface {
    LogUpdate(ID int, typ string) error
}

func Configure(b *tb.Bot) (*tb.Bot, error) {
    b.Handle("/hello", func(m *tb.Message) {
        if _, err := b.Send(m.Sender, "Hello world!"); err != nil {
            log.WithFields(log.Fields{
                "context": "BOT",
                "action":  "SEND",
            }).Error(err)
        }
    })
    return b, nil
}
