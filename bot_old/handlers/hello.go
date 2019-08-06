package handlers

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) HandleHello(update tgbotapi.Update) error {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
    msg.ReplyToMessageID = update.Message.MessageID

    _, err := b.Send(msg)
    return err
}
