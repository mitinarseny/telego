package handlers

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Handler) HandleHello(m *tb.Message) error {
    _, err := b.Bot.Send(m.Sender, "Hello, world!")
    return err
}
