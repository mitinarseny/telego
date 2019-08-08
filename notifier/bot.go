package notifier

import (
    "fmt"

    "github.com/mitinarseny/telego/helpers"
    log "github.com/sirupsen/logrus"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
    *tb.Bot
    Chat *tb.Chat
}

func (b *Bot) String() string {
    return b.Me.Username
}

func (b *Bot) OnNotify(n *Notification) {
    _, err := b.Bot.Send(b.Chat,
        fmt.Sprintf("@%s*: %s*", helpers.EscMd(n.About), n.What),
        &tb.SendOptions{
            ParseMode: tb.ModeMarkdown,
        })
    if err != nil {
        log.WithFields(log.Fields{
            "context": "NOTIFIER",
            "action":  "SEND",
        }).Error(err)
        return
    }
    log.WithFields(log.Fields{
        "context":  "NOTIFIER",
        "notifier": b.Me.Username,
        "action":   "SEND",
    }).Info()
}
