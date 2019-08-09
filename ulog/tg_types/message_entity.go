package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type MessageEntityType string

const (
    HashtagMessageEntityType     MessageEntityType = "hashtag"
    CashtagMessageEntityType     MessageEntityType = "cashtag"
    BotCommandMessageEntityType  MessageEntityType = "bot_command"
    URLMessageEntityType         MessageEntityType = "url"
    EmailMessageEntityType       MessageEntityType = "email"
    PhoneNumberMessageEntityType MessageEntityType = "phone_number"
    BoldMessageEntityType        MessageEntityType = "bold"
    ItalicMessageEntityType      MessageEntityType = "italic"
    CodeEntityType               MessageEntityType = "code"
    PreMessageEntityType         MessageEntityType = "pre"
    TextLinkMessageEntityType    MessageEntityType = "text_link"
    TextMentionMessageEntityType MessageEntityType = "text_mention"
)

type MessageEntity struct {
    Type   MessageEntityType
    Offset int
    Length int
    URL    *string
    User   *User
}

func fromTelebotMessageEntity(e *tb.MessageEntity) *MessageEntity {
    en := new(MessageEntity)
    en.Type = MessageEntityType(e.Type)
    en.Offset = e.Offset
    en.Length = e.Length
    if e.URL != "" {
        en.URL = &e.URL
    }
    if e.User != nil {
        en.User = fromTelebotUser(e.User)
    }
    return en
}
