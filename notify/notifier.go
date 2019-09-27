package notify

type Notification interface {
    Msg() string
}

type Notifier interface {
    Notify(dest string, n Notification) error
}

type AutoNotifier interface {
    Notify(Notification) error
}
