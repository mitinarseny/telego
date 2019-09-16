package handlers

import (
    "errors"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Start struct {
    B *tb.Bot
}

func (h *Start) HandleMsg(m *tb.Message) error {
    return errors.New("handleStart is not implemented yet")
}
