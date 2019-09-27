package handlers

import (
    "github.com/mitinarseny/telego/admins"
    "github.com/mitinarseny/telego/log"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

type CustomizeNotificationsStorage struct {
    Admins admins.AdminsRepo
}

type CustomizeNotifications struct {
    log.UnsafeErrorLogger
    B       *tb.Bot
    Storage *CustomizeNotificationsStorage
}

func (h *CustomizeNotifications) HandleCallback(c *tb.Callback) error {
    _ = h.B.Respond(c, &tb.CallbackResponse{
        Text: "CustomizeNotifications is not implemented yet",
    })
    return errors.New("CustomizeNotifications is not implemented yet")
}
