package tg

type CallbackQuery struct {
    ID              string   `bson:"_id,omitempty"`
    From            User     `bson:"from,omitempty"`
    Message         *Message `bson:"message,omitempty"`
    InlineMessageID *string  `bson:"inline_message_id,omitempty"`
    ChatInstance    string   `bson:"chat_instance,omitempty"`
    Data            *string  `bson:"data,omitempty"`
    GameShortName   *string  `bson:"game_short_name,omitempty"`
}
