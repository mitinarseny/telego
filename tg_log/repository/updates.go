package repository

import (
    "context"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Update struct {
    UpdateID           int64
    Message            *Message
    EditedMessage      *Message
    ChannelPost        *Message
    EditedChannelPost  *Message
    InlineQuery        *InlineQuery
    ChosenInlineResult *ChosenInlineResult
    CallbackQuery      *CallbackQuery
    ShippingQuery      *ShippingQuery
    PreCheckoutQuery   *PreCheckoutQuery
    Poll               *Poll
}

type UpdatesRepository interface {
    Create(ctx context.Context, updates ...*Update) error
}

func FromTelebotUpdate(u *tb.Update) *Update {
    upd := new(Update)
    upd.UpdateID = int64(u.ID)
    if u.Message != nil {
        upd.Message = fromTelebotMessage(u.Message)
    }
    if u.EditedMessage != nil {
        upd.EditedMessage = fromTelebotMessage(u.EditedMessage)
    }
    if u.ChannelPost != nil {
        upd.ChannelPost = fromTelebotMessage(u.ChannelPost)
    }
    if u.EditedChannelPost != nil {
        upd.EditedChannelPost = fromTelebotMessage(u.EditedChannelPost)
    }
    if u.Query != nil {
        upd.InlineQuery = fromTelebotQuery(u.Query)
    }
    if u.ChosenInlineResult != nil {
        upd.ChosenInlineResult = fromTelebotChosenInlineResult(u.ChosenInlineResult)
    }
    if u.PreCheckoutQuery != nil {
        upd.PreCheckoutQuery = fromPreCheckoutQuery(u.PreCheckoutQuery)
    }
    return upd
}
