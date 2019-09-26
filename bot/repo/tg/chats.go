package tg

import (
    "context"

    "github.com/mitinarseny/telego/bot/repo"
)

type ChatType string

const (
    PrivateChatType    ChatType = "private"
    GroupChatType      ChatType = "group"
    SupergroupChatType ChatType = "supergroup"
    ChannelChatType    ChatType = "channel"
)

type Chat struct {
    repo.BaseModel `bson:",inline"`

    ID        int64    `bson:"_id,omitempty"`
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
