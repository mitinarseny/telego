package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleHello(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(msg)
	return err
}
