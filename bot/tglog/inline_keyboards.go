package tglog

type InlineKeyboardMarkup struct {
    InlineKeyboard [][]InlineKeyboardButton `bson:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
    Text                         string    `bson:"text,omitempty"`
    URL                          *string   `bson:"url,omitempty"`
    LoginURL                     *LoginURL `bson:"login_url,omitempty"`
    CallbackData                 *string   `bson:"callback_data,omitempty"`
    SwitchInlineQuery            *string   `bson:"switch_inline_query,omitempty"`
    SwitchInlineQueryCurrentChat *string   `bson:"switch_inline_query_current_chat,omitempty"`
    CallbackGame                 *Game     `bson:"callback_game,omitempty"`
    Pay                          *bool     `bson:"pay,omitempty"`
}

type LoginURL struct {
    URL                string  `bson:"url,omitempty"`
    ForwardText        *string `bson:"forward_text,omitempty"`
    BotUsername        *string `bson:"bot_username,omitempty"`
    RequestWriteAccess *bool   `bson:"request_write_access,omitempty"`
}
