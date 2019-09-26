package handlers

import (
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Logger interface {
    Log(error)
}

func MsgWithLog(l Logger, h MsgHandler) func(*tb.Message) {
    return func(m *tb.Message) {
        if err := h.HandleMsg(m); err != nil {
            log.Error(err) // log with ErrorLogger
        }
    }
}

func CallbackWithLog(l Logger, h CallbackHandler) func(*tb.Callback) {
    return func(c *tb.Callback) {
        if err := h.HandleCallback(c); err != nil {
            log.Error(err) // log with ErrorLogger
        }
    }
}
