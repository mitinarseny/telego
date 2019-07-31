package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func notify(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	_, err := bot.Send(tgbotapi.NewMessage(chatID, text))
	if err != nil {
		log.Printf("Failed to notify with @%s", bot.Self.UserName)
	}
	return err
}

func botStatusText(botUserName, status string) string {
	return fmt.Sprintf("@%s: %s", botUserName, status)
}
