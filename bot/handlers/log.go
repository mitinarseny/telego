package handlers

import (
    "github.com/mitinarseny/telego/log"
    tb "gopkg.in/tucnak/telebot.v2"
)

func MsgWithLog(l log.UnsafeInfoErrorLogger, h MsgHandler) func(*tb.Message) {
    return func(m *tb.Message) {
        if err := h.HandleMsg(m); err != nil {
            l.Error(err)
        }
    }
}

func CallbackWithLog(l log.UnsafeInfoErrorLogger, h CallbackHandler) func(*tb.Callback) {
    return func(c *tb.Callback) {
        if err := h.HandleCallback(c); err != nil {
            l.Error(err)
        }
    }
}
