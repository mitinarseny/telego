package notifier

import (
    "fmt"

    "github.com/mitinarseny/telego/helpers"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
    *tb.Bot
    Chat *tb.Chat
}

func (b *Bot) Notify(about, text string) error {
    _, err := b.Bot.Send(b.Chat,
        fmt.Sprintf("@%s*: %s*", helpers.EscMd(about), text),
        tb.SendOptions{
            ParseMode: tb.ModeMarkdown,
        })
    return err
}

func (b *Bot) Notifyf(about, format string, args ...interface{}) error {
    _, err := b.Bot.Send(b.Chat,
        fmt.Sprintf("@%s*: %s*", helpers.EscMd(about), fmt.Sprintf(format, args)),
        tb.SendOptions{
            ParseMode: tb.ModeMarkdown,
        })
    return err
}

func (b *Bot) NotifyError(about string, err error) error {
    _, er := b.Bot.Send(b.Chat,
        fmt.Sprintf("@%s*: %s*", helpers.EscMd(about), err))
    return er
}

func (b *Bot) NotifyUp(about string) error {
    return b.Notifyf(about, upMessage)
}

func (b *Bot) NotifyDown(about string) error {
    return b.Notifyf(about, downMessage)
}
