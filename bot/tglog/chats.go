package tglog

import (
    "context"

    tb "gopkg.in/tucnak/telebot.v2"
)

type ChatType string

const (
    PrivateChatType    ChatType = "private"
    GroupChatType      ChatType = "group"
    SupergroupChatType ChatType = "supergroup"
    ChannelChatType    ChatType = "channel"
)

type Chat struct {
    ID        int64    `bson:"id,omitempty"`
    Type      ChatType `bson:"type,omitempty"`
    Title     *string  `bson:"title,omitempty"`
    Username  *string  `bson:"username,omitempty"`
    FirstName *string  `bson:"first_name,omitempty"`
    LastName  *string  `bson:"last_name,omitempty"`
}

type ChatsRepo interface {
    Create(ctx context.Context, chats ...*Chat) ([]*Chat, error)
    CreateIfNotExist(ctx context.Context, chats ...*Chat) ([]*Chat, error)
}

func fromTelebotChat(c *tb.Chat) *Chat {
    return &Chat{
        ID:        c.ID,
        Type:      ChatType(c.Type),
        Title:     &c.Title,
        Username:  &c.Username,
        FirstName: &c.FirstName,
        LastName:  &c.LastName,
    }
}
