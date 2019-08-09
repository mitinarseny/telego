package tg_types

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
