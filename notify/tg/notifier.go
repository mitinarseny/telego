package tg

import (
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/notify"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

const (
    notificationFormat = `ℹ️ *%s*

_To change notification settings send_ %s`
)

type Notifier struct {
    tg                          *tb.Bot
    notificationSettingsCommand string
}

func NewNotifier(tg *tb.Bot, notificationSettingsCommand string) *Notifier {
    return &Notifier{
        tg:                          tg,
        notificationSettingsCommand: notificationSettingsCommand,
    }
}

func (n *Notifier) Notify(dest string, nt notify.Notification) error {
    text := fmt.Sprintf(notificationFormat, nt.Msg(), n.notificationSettingsCommand) // TODO: nicer message style
    userID, err := strconv.Atoi(dest)
    if err != nil {
        return errors.Wrapf(err, "can not parse userID from %q", dest)
    }
    _, err = n.tg.Send(&tb.User{ID: userID}, text, &tb.SendOptions{
        ParseMode: tb.ModeMarkdown,
    })
    return err
}
