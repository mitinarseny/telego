package tg

import (
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/notify"
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

func (n *Notifier) Notify(dest string, nt notify.Notification) error {
    text := fmt.Sprintf("*%s*", nt.Msg()) // TODO: nicer message style
    userID, err := strconv.Atoi(dest)
    if err != nil {
        return errors.Wrapf(err, "can not parse userID from %q", dest)
    }
    _, err = n.tg.Send(&tb.User{ID: userID}, text, &tb.SendOptions{
        ParseMode: tb.ModeMarkdown,
        // TODO: add buttons to customise notifications settings for current user
    })
    return err
}
