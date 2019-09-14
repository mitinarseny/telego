package repo

type CallbackQuery struct {
    ID              string
    From            User
    Message         *Message
    InlineMessageID *string
    ChatInstance    string
    Data            *string
    GameShortName   *string
}
