package filters

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgFilter interface {
    FilterMsg(*tb.Message) (bool, error)
}

type CallbackFilter interface {
    FilterCallback(*tb.Callback) (bool, error)
}
