package tg

import (
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/bot/notifier"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Notifier struct {
    tg *tb.Bot
}

func NewNotifier(tg *tb.Bot) *Notifier {
    return &Notifier{
        tg: tg,
    }
}

func (n *Notifier) Notify(dest string, notification *notifier.Notification) error {
    text := fmt.Sprintf("*%s*", notification.Msg) // TODO: nicer message style
    userID, err := strconv.Atoi(dest)
    if err != nil {
        return errors.Wrapf(err, "can not convert %s to int", dest)
    }
    _, err = n.tg.Send(&tb.User{ID: userID}, text, &tb.SendOptions{
        ParseMode: tb.ModeMarkdown,
        // TODO: add buttons to customise notification settings for current user
    })
    return err
}
