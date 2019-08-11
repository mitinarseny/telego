package repository

import (
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
    ID        int64
    Type      ChatType
    Title     *string
    Username  *string
    FirstName *string
    LastName  *string
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
