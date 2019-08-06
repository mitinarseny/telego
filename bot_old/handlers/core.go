package handlers

import (
    "github.com/go-telegram-bot-api/telegram-bot-api"
    "github.com/mitinarseny/telego/bot_old/ch_log"
    "github.com/mitinarseny/telego/bot_old/ulog"
)

type Bot struct {
    *tgbotapi.BotAPI
}

func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel, ul ch_log.UpdatesLogger, errCh chan<- error) error {
    go func() {
        ul.LogUpdates(updates)
    }()
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
