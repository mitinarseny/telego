package bot_old

import (
    "fmt"

    log "github.com/sirupsen/logrus"

    "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
    UpMessage   = "UP ❇️"
    DownMessage = "DOWN ❗️"
)

func notify(bot *tgbotapi.BotAPI, chatID int64, text string) error {
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = tgbotapi.ModeMarkdown

    _, err := bot.Send(msg)
    if err != nil {
        log.WithField("notifier", bot.Self.UserName).Warn("Failed to notify")
    }
    return err
}

func botStatusText(botUserName, status string) string {
    return fmt.Sprintf("@%s *is %s*", escMd(botUserName), status)
}

func notifyUp(notifier *tgbotapi.BotAPI, chatID int64, botName string) error {
    return notify(notifier, chatID, botStatusText(botName, UpMessage))
}

func notifyDown(notifier *tgbotapi.BotAPI, chatID int64, botName string) error {
    return notify(notifier, chatID, botStatusText(botName, DownMessage))
}
