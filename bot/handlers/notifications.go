package handlers

import (
    "github.com/mitinarseny/telego/admins"
    "github.com/mitinarseny/telego/log"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

type NotificationsStorage struct {
    Admins admins.AdminsRepo
}

type notifications struct {
    log.UnsafeErrorLogger
    tg      *tb.Bot
    storage *NotificationsStorage
}

type NotificationsSettings struct {
    Logger  log.UnsafeInfoErrorLogger
    Tg      *tb.Bot
    Storage *NotificationsStorage
}

func NewNotifications(pref *NotificationsSettings) *notifications {
    return &notifications{
        UnsafeErrorLogger: pref.Logger,
        tg:                pref.Tg,
        storage:           pref.Storage,
    }
}

func (h *notifications) HandleMsg(m *tb.Message) error {
    return errors.New("notifications is not implemented yet")
}
