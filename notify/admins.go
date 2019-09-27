package notify

import (
    "context"
    "strconv"

    "github.com/mitinarseny/telego/admins"
    "github.com/mitinarseny/telego/log"
    "github.com/pkg/errors"
)

type AdminsNotifier struct {
    ErrorLogger log.UnsafeInfoErrorLogger
    Admins      admins.AdminsRepo
    Notifiers   map[admins.NotifierType]Notifier
}

func (n *AdminsNotifier) Notify(nt Notification) error {
    var (
        adms                 []*admins.Admin
        extractNotifierTypes func(*admins.Admin) []admins.NotifierType
        err                  error
    )
    switch nt.(type) {
    case StatusNotification:
        adms, err = n.Admins.GetAllShouldBeNotifiedAbout(context.Background(), admins.StatusNotificationType)
        extractNotifierTypes = func(a *admins.Admin) []admins.NotifierType {
            return a.Notifications.Status
        }
    default:
        return errors.Errorf("unsupported notification type %T", nt)
    }
    if err != nil {
        return err
    }
    for _, a := range adms {
        types := extractNotifierTypes(a)
        for _, typ := range types {
            ntf, found := n.Notifiers[typ]
            if !found {
                n.ErrorLogger.Error(errors.Errorf("unknown nn type %q", typ))
            }
            if err := ntf.Notify(strconv.FormatInt(a.ID, 10), nt); err != nil {
                n.ErrorLogger.Error(err)
            }
        }
    }
    return nil
}

func (n *AdminsNotifier) NotifyStatus(s Status) error {
    return n.Notify(StatusNotification(s))
}

func (n *AdminsNotifier) NotifyStatusUp() error {
    return n.NotifyStatus(StatusUp)
}

func (n *AdminsNotifier) NotifyStatusDown() error {
    return n.NotifyStatus(StatusDown)
}
