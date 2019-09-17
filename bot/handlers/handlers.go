package handlers

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgHandler interface {
    HandleMsg(*tb.Message) error
}

type CallbackHandler interface {
    HandleCallback(*tb.Callback) error
}
