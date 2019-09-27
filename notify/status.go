package notify

import "fmt"

const (
    StatusUp   Status = "UP"
    StatusDown Status = "DOWN"
)

type Status string

type StatusNotification Status

func (n StatusNotification) Msg() string {
    return fmt.Sprintf("Staus: %s", n)
}
