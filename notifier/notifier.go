package notifier

import (
    "context"
    "sync"

    log "github.com/sirupsen/logrus"
)

const (
    upMessage   = "UP ❇️"
    downMessage = "DOWN ❗️"
)

type Notification struct {
    About string
    What  interface{}
}

// Observer is someone who handles notifications
type Observer interface {
    String() string
    OnNotify(*Notification)
}

// Notifier is someone who is responsible for delivering notifications to Observers
type Notifier interface {
    // Register register an Observer for further handling events
    Register(Observer)
    // Deregister deletes an Observer from its observers
    Deregister(Observer)
    // Start starts delivering Notifications to observers
    Start(ctx context.Context)
    // // Stop stops delivering Notifications to observers
    Stop()
    // Notify delivers Notification to registered Observers
    Notify(*Notification)
}

// StatusNotifier is someone who is responsible for delivering Up/Down notifications to Observers
type StatusNotifier interface {
    Notifier
    // NotifyUP delivers an UP notification to its Observers
    NotifyUp(about string)
    // NotifyDown delivers a Down notification to its Observers
    NotifyDown(about string)
}

type BaseNotifier struct {
    observers map[Observer]struct{}
    mu        sync.Mutex
    ch        chan *Notification
    wg        sync.WaitGroup
}

func NewBaseNotifier() *BaseNotifier {
    return &BaseNotifier{
        observers: make(map[Observer]struct{}),
    }
}

func (n *BaseNotifier) Register(o Observer) {
    n.mu.Lock()
    defer n.mu.Unlock()
    log.WithFields(log.Fields{
        "context": "NOTIFIER",
        "action": "REGISTER",
        "notifier": o,
    }).Info()
    n.observers[o] = struct{}{}
}

func (n *BaseNotifier) Deregister(o Observer) {
    n.mu.Lock()
    defer n.mu.Unlock()

    log.WithFields(log.Fields{
        "context": "NOTIFIER",
        "action": "DEREGISTER",
        "notifier": o,
    }).Info()
    delete(n.observers, o)
}

func (n *BaseNotifier) Start(ctx context.Context) {
    n.ch = make(chan *Notification)

    go func() {
        for {
            select {
            case nt, more := <-n.ch:
                if !more {
                    return
                }
                var wg sync.WaitGroup
                wg.Add(len(n.observers))
                for o := range n.observers {
                    go func() {
                        defer wg.Done()
                        o.OnNotify(nt)
                    }()
                }
                go func() {
                    defer n.wg.Done()
                    wg.Wait()
                }()
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (n *BaseNotifier) Stop() {
    defer close(n.ch)
    n.wg.Wait()
}

func (n *BaseNotifier) Notify(nt *Notification) {
    n.wg.Add(1)
    go func() {
        n.ch <- nt
    }()
}

func (n *BaseNotifier) NotifyUp(about string) {
    n.Notify(&Notification{
        About: about,
        What:  upMessage,
    })
}

func (n *BaseNotifier) NotifyDown(about string) {
    n.Notify(&Notification{
        About: about,
        What:  downMessage,
    })
}
