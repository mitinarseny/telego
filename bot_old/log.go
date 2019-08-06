package bot_old

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type UpdatesLogger interface {
    LogUpdates(tgbotapi.UpdatesChannel)
}
