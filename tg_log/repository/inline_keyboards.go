package repository

type InlineKeyboardMarkup struct {
    InlineKeyboard [][]InlineKeyboardButton
}

type InlineKeyboardButton struct {
    Text                         string
    URL                          *string
    LoginURL                     *LoginURL
    CallbackData                 *string
    SwitchInlineQuery            *string
    SwitchInlineQueryCurrentChat *string
    CallbackGame                 *Game
    Pay                          *bool
}

type LoginURL struct {
    URL                string
    ForwardText        *string
    BotUsername        *string
    RequestWriteAccess *bool
}
