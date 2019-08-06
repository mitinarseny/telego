package handlers

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
    unsupportedAnswer = "I was not expecting to get such a message. Please, try again."
)

func (b *Bot) HandleUnsupported(update tgbotapi.Update) error {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, unsupportedAnswer)
    msg.ParseMode = tgbotapi.ModeMarkdown
    msg.ReplyToMessageID = update.Message.MessageID

    _, err := b.Send(msg)
    return err
}
