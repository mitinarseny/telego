package admins

import (
    "context"
    "fmt"

    "github.com/mitinarseny/telego/bot/notifier"
    "github.com/mitinarseny/telego/bot/repo/administration"
)

type Status string

const (
    StatusUp   Status = "UP"
    StatusDown Status = "DOWN"
)

func (n *Notifier) NotifyStatus(s Status) error {
    admins, err := n.Admins.GetAllShouldBeNotifiedAbout(context.Background(), administration.StatusNotificationType)
    if err != nil {
        return err
    }
    for _, a := range admins {
        if err := n.notify(a,
            notifier.NewNotification(fmt.Sprintf("I am %s", s)),
            a.Notifications.Status...); err != nil {
            n.ErrorLogger.Log(err)
        }
    }
    return nil
}
