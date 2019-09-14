package repo

import (
    "context"
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Update struct {
    CreatedAt *time.Time `bson:"_created_at,omitempty"`

    UpdateID           int64               `bson:"_id,omitempty"`
    Message            *Message            `bson:"message,omitempty"`
    EditedMessage      *Message            `bson:"edited_message,omitempty"`
    ChannelPost        *Message            `bson:"channel_post,omitempty"`
    EditedChannelPost  *Message            `bson:"edited_channel_post,omitempty"`
    InlineQuery        *InlineQuery        `bson:"inline_query,omitempty"`
    ChosenInlineResult *ChosenInlineResult `bson:"chosen_inline_result,omitempty"`
    CallbackQuery      *CallbackQuery      `bson:"callback_query,omitempty"`
    ShippingQuery      *ShippingQuery      `bson:"shipping_query,omitempty"`
    PreCheckoutQuery   *PreCheckoutQuery   `bson:"pre_checkout_query,omitempty"`
    Poll               *Poll               `bson:"poll,omitempty"`
}

func (u *Update) ID() interface{} {
    return u.UpdateID
}

type UpdatesRepo interface {
    Create(ctx context.Context, updates []Update) ([]Update, error)
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
