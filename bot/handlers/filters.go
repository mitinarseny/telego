package handlers

import (
    "github.com/mitinarseny/telego/bot/filters"
    tb "gopkg.in/tucnak/telebot.v2"
)

type withMsgFilters struct {
    handler MsgHandler
    filters []filters.MsgFilter
}

func WithMsgFilters(h MsgHandler, fs ...filters.MsgFilter) *withMsgFilters {
    return &withMsgFilters{
        handler: h,
        filters: fs,
    }
}

func (h *withMsgFilters) HandleMsg(m *tb.Message) error {
    for _, f := range h.filters {
        if passed, err := f.FilterMsg(m); err != nil {
            return err
        } else if !passed {
            return nil
        }
    }
    return h.handler.HandleMsg(m)
}

type withCallbackFilters struct {
    handler CallbackHandler
    filters []filters.CallbackFilter
}

func WithCallbackFilters(h CallbackHandler, fs ...filters.CallbackFilter) *withCallbackFilters {
    return &withCallbackFilters{
        handler: h,
        filters: fs,
    }
}

func (h *withCallbackFilters) HandleCallback(c *tb.Callback) error {
    for _, f := range h.filters {
        if passed, err := f.FilterCallback(c); err != nil {
            return err
        } else if !passed {
            return nil
        }
    }
    return h.handler.HandleCallback(c)
}
