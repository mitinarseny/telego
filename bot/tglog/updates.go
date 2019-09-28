package tglog

import (
    "context"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Update struct {
    UpdateID           int64               `bson:"update_id,omitempty"`
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

func fromTelebotUpdate(u *tb.Update) *Update {
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
    if u.Callback != nil {
        upd.CallbackQuery = fromTelebotCallback(u.Callback)
    }
    if u.PreCheckoutQuery != nil {
        upd.PreCheckoutQuery = fromPreCheckoutQuery(u.PreCheckoutQuery)
    }
    return upd
}

type UpdatesRepo interface {
    Create(ctx context.Context, updates ...*Update) ([]*Update, error)
    CreateIfNotExist(ctx context.Context, updates ...*Update) ([]*Update, error)
}

type Logger struct {
    r UpdatesRepo
}

func NewUpdatesLogger(r UpdatesRepo) *Logger {
    return &Logger{
        r: r,
    }
}

func (l *Logger) LogUpdates(updates ...*tb.Update) error {
    models := make([]*Update, 0, len(updates))
    for _, u := range updates {
        models = append(models, fromTelebotUpdate(u))
    }
    _, err := l.r.Create(context.Background(), models...)
    return err
}
