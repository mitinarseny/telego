package filters

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgFilter interface {
    Filter(*tb.Message) (bool, error)
}
