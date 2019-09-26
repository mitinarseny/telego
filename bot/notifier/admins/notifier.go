package admins

import (
    "strconv"

    errlog "github.com/mitinarseny/telego/bot/logging/errors"
    "github.com/mitinarseny/telego/bot/notifier"
    "github.com/mitinarseny/telego/bot/notifier/tg"
    "github.com/mitinarseny/telego/bot/repo/administration"
    "github.com/pkg/errors"
)

type Notifier struct {
    ErrorLogger errlog.Logger
    Admins      administration.AdminsRepo
    Notifiers   map[administration.NotifierType]notifier.Notifier
}

func (n *Notifier) notify(admin *administration.Admin,
    notification *notifier.Notification,
    destinations ...administration.NotifierType) error {
    for _, typ := range destinations {
        nn, found := n.Notifiers[typ]
        if !found {
            return errors.Errorf("unknown nn type %q", typ)
        }
        var who string
        switch nn.(type) {
        case *tg.Notifier:
            who = strconv.FormatInt(admin.ID, 10)
        default:
            return errors.Errorf("unsupported nn type %T", nn)
        }
        if err := nn.Notify(who, notification); err != nil {
            return err
        }
    }
    return nil
}
