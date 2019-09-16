package handlers

import (
    "github.com/mitinarseny/telego/bot/filters"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgHandler interface {
    HandleMsg(*tb.Message) error
}

type CallbackHandler interface {
    HandleCallback(*tb.Callback) error
}

type Logger interface {
    Log(error)
}

func MsgWithLog(l Logger, h MsgHandler) func(*tb.Message) {
    return func(m *tb.Message) {
        if err := h.HandleMsg(m); err != nil {
            l.Log(err)
        }
    }
}

func CallbackWithLog(l Logger, h CallbackHandler) func(*tb.Callback) {
    return func(c *tb.Callback) {
        if err := h.HandleCallback(c); err != nil {
            l.Log(err)
        }
    }
}

type withFilters struct {
    handler MsgHandler
    filters []filters.MsgFilter
}

func MsgWithFilters(h MsgHandler, filters ...filters.MsgFilter) *withFilters {
    return &withFilters{
        handler: h,
        filters: filters,
    }
}

func (h *withFilters) HandleMsg(m *tb.Message) error {
    for _, f := range h.filters {
        passed, err := f.Filter(m)
        if err != nil {
            return err
        }
        if !passed {
            return nil
        }
    }
    return h.handler.HandleMsg(m)
}
