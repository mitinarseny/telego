package bot

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mitinarseny/telego/bot/handlers"
)

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	// main logic here
	switch {
	case update.Message != nil:
		switch {
		case update.Message.Command() == "hello":
			return handlers.HandleHello(bot, update)
		}

	}
	return nil
}

func Start(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account: @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if err := handleUpdate(bot, update); err != nil {
			return err
		}
	}

	return nil
}
