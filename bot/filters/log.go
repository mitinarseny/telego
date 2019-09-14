package filters

import (
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

func (f *Filters) WithLog(h MessageHandler) func(*tb.Message) {
    return func(m *tb.Message) {
        if err := h(m); err != nil {
            log.WithFields(log.Fields{
                "context": "BOT",
                "handler": h,
            }).Error(err)
        }
    }
}
