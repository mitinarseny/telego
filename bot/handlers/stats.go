package handlers

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

func (h *Handler) HandleStats(m *tb.Message) error {
    _, err := h.Bot.Send(m.Sender, "Hello, world!", &tb.SendOptions{
        ReplyTo: m,
    })
    return err
}
