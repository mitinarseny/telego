package chattools

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
    b *tb.Bot
}

func WithBot(b *tb.Bot) *bot {
    return &bot{
        b: b,
    }
}
