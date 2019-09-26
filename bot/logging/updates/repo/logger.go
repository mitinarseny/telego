package repo

import (
    "context"

    "github.com/mitinarseny/telego/bot/repo/tg"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Logger struct {
    r tg.UpdatesRepo
}

func NewUpdatesLogger(r tg.UpdatesRepo) *Logger {
    return &Logger{
        r: r,
    }
}

func (l *Logger) LogUpdates(updates ...*tb.Update) error {
    models := make([]*tg.Update, 0, len(updates))
    for _, u := range updates {
        models = append(models, fromTelebotUpdate(u))
    }
    _, err := l.r.Create(context.Background(), models...)
    return err
}
