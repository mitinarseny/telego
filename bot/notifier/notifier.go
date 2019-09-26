package notifier

type Notification struct {
    Msg string
}

func NewNotification(msg string) *Notification {
    return &Notification{
        Msg: msg,
    }
}

type Notifier interface {
    Notify(dest string, n *Notification) error
}


