package bot

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mitinarseny/telego/bot/handlers"
)

const (
	notifyParseMode = "Markdown"
)

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	// main logic here
	switch {
	case update.Message != nil:
		switch {
		case update.Message.Command() == "hello":
			return handlers.HandleHello(bot, update)
		default:
			return handlers.HandleUnsupported(bot, update)
		}
	default:
		return handlers.HandleUnsupported(bot, update)
	}
}

func Run(token, notifierToken string, notifyChatID int64, debug bool) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	bot.Debug = debug
	log.Printf("Authorized on account: @%s", bot.Self.UserName)

	if notifierToken != "" {
		notifier, err := tgbotapi.NewBotAPI(notifierToken)
		if err != nil {
			return err
		}
		log.Printf("Notifier: @%s", notifier.Self.UserName)
		defer func() {
			_ = notify(notifier, notifyChatID, botStatusText(bot.Self.UserName, "down ❗️"))
		}()
		_ = notify(notifier, notifyChatID, botStatusText(bot.Self.UserName, "up ❇️"))
	}

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
