package bot

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/mitinarseny/telego/bot/filters"
    "github.com/mitinarseny/telego/bot/handlers"
    errlog "github.com/mitinarseny/telego/bot/logging/errors"
    "github.com/mitinarseny/telego/bot/logging/updates"
    "github.com/mitinarseny/telego/bot/notifier"
    "github.com/mitinarseny/telego/bot/notifier/admins"
    "github.com/mitinarseny/telego/bot/notifier/tg"
    "github.com/mitinarseny/telego/bot/repo/administration"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgHandlerFunc func(*tb.Message) error

type MsgFilterFunc func(*tb.Message) (bool, error)

type Storage struct {
    Admins administration.AdminsRepo
    Roles  administration.RolesRepo
}

type Logger interface {
    Log(error)
}

type Bot struct {
    tg *tb.Bot

    s *Storage

    errLogger    errlog.Logger
    updateLogger updates.Logger
    Notifier     *admins.Notifier
    superuserID  int64
}

type Settings struct {
    Token        string
    LastUpdateID int
    Storage      *Storage
    UpdateLogger updates.Logger
    ErrorLogger  errlog.Logger
    SuperuserID  int64
}

func New(s *Settings) (*Bot, error) {
    b := Bot{
        s:            s.Storage,
        errLogger:    s.ErrorLogger,
        updateLogger: s.UpdateLogger,
        superuserID:  s.SuperuserID,
    }
    bot, err := tb.NewBot(tb.Settings{
        Token: s.Token,
        Poller: tb.NewMiddlewarePoller(&tb.LongPoller{
            Timeout:      60 * time.Second,
            LastUpdateID: s.LastUpdateID,
        }, func(u *tb.Update) bool {
            go func() {
                if err := b.updateLogger.LogUpdates(u); err != nil {
                    b.errLogger.Log(err)
                }
            }()
            return true
        }),
        Client: &http.Client{
            Timeout: 0,
        },
        Reporter: func(err error) {
            b.errLogger.Log(err)
        },
    })
    if err != nil {
        return nil, err
    }
    b.tg = bot
    b.Notifier = &admins.Notifier{
        ErrorLogger: s.ErrorLogger,
        Admins:      s.Storage.Admins,
        Notifiers: map[administration.NotifierType]notifier.Notifier{
            administration.TelegramNotificationDest: tg.NewNotifier(b.tg),
        },
    }
    if err := b.setSuperuser(s.SuperuserID); err != nil {
        return nil, err
    }
    b.setupHandlers()
    return &b, nil
}

func (b *Bot) setSuperuser(userID int64) error {
    userIDStr := strconv.FormatInt(userID, 10)
    if _, err := b.tg.ChatByID(userIDStr); err != nil {
        return errors.Wrapf(err,
            "can't get chat with superuser %q, check that superuser has a conversation with bot",
            userIDStr)
    }
    _, err := b.s.Admins.CreateIfNotExists(context.Background(), &administration.Admin{
        ID:   userID,
        Role: administration.SuperuserRole,
        Notifications: administration.NotificationsPreferences{
            Status: []administration.NotifierType{
                administration.TelegramNotificationDest,
            },
        },
    })
    return err
}

func (b *Bot) setupHandlers() {
    b.tg.Handle("/start", handlers.MsgWithLog(b.errLogger, &handlers.Start{
        B: b.tg,
    }))
    b.tg.Handle("/admins", handlers.MsgWithLog(b.errLogger, handlers.WithMsgFilters(&handlers.Admins{
        Logger: b.errLogger,
        B:      b.tg,
        Storage: &handlers.AdminsStorage{
            Admins: b.s.Admins,
            Roles:  b.s.Roles,
        },
    }, filters.WithSender().IsAdminWithScopes(b.s.Admins, administration.AdminsReadScope))))
}

func (b *Bot) Start() {
    if err := b.Notifier.NotifyStatus(admins.StatusUp); err != nil {
        b.errLogger.Log(err)
    }
    b.tg.Start()
}

func (b *Bot) Stop() {
    defer func() {
        if err := b.Notifier.NotifyStatus(admins.StatusDown); err != nil {
            b.errLogger.Log(err)
        }
    }()
    b.tg.Stop()
}
