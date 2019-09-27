package tg

import (
    "fmt"
    "strconv"

    "github.com/mitinarseny/telego/notify"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

const (
    notificationFormat = "ℹ️ *%s*"
)

type Notifier struct {
    tg                *tb.Bot
    customizeEndpoint string
}

func NewNotifier(tg *tb.Bot, customizeEndpoint string) *Notifier {
    return &Notifier{
        tg:                tg,
        customizeEndpoint: customizeEndpoint,
    }
}

func (n *Notifier) Notify(dest string, nt notify.Notification) error {
    text := fmt.Sprintf(notificationFormat, nt.Msg()) // TODO: nicer message style
    userID, err := strconv.Atoi(dest)
    if err != nil {
        return errors.Wrapf(err, "can not parse userID from %q", dest)
    }
    _, err = n.tg.Send(&tb.User{ID: userID}, text, &tb.SendOptions{
        ParseMode: tb.ModeMarkdown,
    }, &tb.ReplyMarkup{
        InlineKeyboard: [][]tb.InlineButton{{{
            Unique: n.customizeEndpoint,
            Text:   "Customize Notifications",
            Data:   "", // TODO: add data to identify or with sessions ???
        }}},
    })
    return err
}
