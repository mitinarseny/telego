package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type CallbackQuery struct {
    ID              string   `bson:"id,omitempty"`
    From            User     `bson:"from,omitempty"`
    Message         *Message `bson:"message,omitempty"`
    InlineMessageID *string  `bson:"inline_message_id,omitempty"`
    ChatInstance    string   `bson:"chat_instance,omitempty"`
    Data            *string  `bson:"data,omitempty"`
    GameShortName   *string  `bson:"game_short_name,omitempty"`
}

func fromTelebotCallback(c *tb.Callback) *CallbackQuery {
    clb := new(CallbackQuery)
    clb.ID = c.ID
    if c.Sender != nil {
        clb.From = *fromTelebotUser(c.Sender)
    }
    if c.Message != nil {
        clb.Message = fromTelebotMessage(c.Message)
    }
    if c.MessageID != "" {
        clb.InlineMessageID = &c.MessageID
    }
    if c.Data != "" {
        clb.Data = &c.Data
    }
    return clb
}
