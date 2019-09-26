package tg

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
    Type   MessageEntityType `bson:"type,omitempty"`
    Offset int               `bson:"offset,omitempty"`
    Length int               `bson:"length,omitempty"`
    URL    *string           `bson:"url,omitempty"`
    User   *User             `bson:"user,omitempty"`
}
