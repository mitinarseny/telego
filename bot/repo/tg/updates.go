package tg

import (
    "context"

    "github.com/mitinarseny/telego/bot/repo"
)

type Update struct {
    repo.BaseModel `bson:",inline"`

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

func (u *Update) From() *User {
    switch {
    case u.Message != nil && u.Message.From != nil:
        return u.Message.From
    case u.EditedMessage != nil && u.EditedMessage.From != nil:
        return u.EditedMessage.From
    case u.ChannelPost != nil && u.ChannelPost.From != nil:
        return u.ChannelPost.From
    case u.EditedChannelPost != nil && u.EditedChannelPost.From != nil:
        return u.EditedChannelPost.From
    case u.ChosenInlineResult != nil:
        return &u.ChosenInlineResult.From
    case u.CallbackQuery != nil:
        return &u.CallbackQuery.From
    case u.ShippingQuery != nil:
        return &u.ShippingQuery.From
    case u.PreCheckoutQuery != nil:
        return &u.PreCheckoutQuery.From
    }
    return nil
}

func (u *Update) Chat() *Chat {
    switch {
    case u.Message != nil && u.Message.Chat != nil:
        return u.Message.Chat
    case u.EditedMessage != nil && u.EditedMessage.Chat != nil:
        return u.EditedMessage.Chat
    case u.ChannelPost != nil && u.ChannelPost.Chat != nil:
        return u.ChannelPost.Chat
    case u.EditedChannelPost != nil && u.EditedChannelPost.Chat != nil:
        return u.EditedChannelPost.Chat
    case u.CallbackQuery != nil && u.CallbackQuery.Message != nil && u.CallbackQuery.Message.Chat != nil:
        return u.CallbackQuery.Message.Chat
    }
    return nil
}

type UpdatesRepo interface {
    Create(ctx context.Context, updates ...*Update) ([]*Update, error)
    CreateIfNotExist(ctx context.Context, updates ...*Update) ([]*Update, error)
}
