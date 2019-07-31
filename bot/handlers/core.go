package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	*tgbotapi.BotAPI
}

func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel, errCh chan <- error) error {
	for update := range updates {
		go func() {
			switch {
			case update.Message != nil:
				switch {
				case update.Message.Command() == "hello":
					if err := b.HandleHello(update); err != nil {
						errCh <- err
					}
				default:
					if err := b.HandleUnsupported(update); err != nil {
						errCh <- err
					}
				}
			default:
				if err := b.HandleUnsupported(update); err != nil {
					errCh <- err
				}
			}
		}()
	}
	return nil
}
