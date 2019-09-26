package updates

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Logger interface {
    LogUpdates(...*tb.Update) error
}
