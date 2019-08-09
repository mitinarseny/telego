package tg_types

type LoginURL struct {
    URL                string
    ForwardText        *string
    BotUsername        *string
    RequestWriteAccess *bool
}
