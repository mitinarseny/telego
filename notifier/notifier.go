package notifier

const (
    upMessage   = "UP ❇️"
    downMessage = "DOWN ❗️"
)

type Notifier interface {
    Notify(about, text string) error
    Notifyf(about, format string, args ...interface{}) error
    NotifyError(about string, err error) error
}

type StatusNotifier interface {
    NotifyUp(about string) error
    NotifyDown(about string) error
}

