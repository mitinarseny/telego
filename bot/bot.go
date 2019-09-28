package bot

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/mitinarseny/telego/admins"
    "github.com/mitinarseny/telego/bot/filters"
    "github.com/mitinarseny/telego/bot/handlers"
    "github.com/mitinarseny/telego/bot/tglog"
    "github.com/mitinarseny/telego/log"
    "github.com/mitinarseny/telego/notify"
    tgnotify "github.com/mitinarseny/telego/notify/tg"
    "github.com/pkg/errors"
    tb "gopkg.in/tucnak/telebot.v2"
)

// Endpoints
const (
    // Commands
    startCommand         = "/start"
    adminsCommand        = "/admins"
    notificationsCommand = "/notifications"
)

type MsgHandlerFunc func(*tb.Message) error

type MsgFilterFunc func(*tb.Message) (bool, error)

type Storage struct {
    Admins admins.AdminsRepo
    Roles  admins.RolesRepo
}

type Bot struct {
    tg             *tb.Bot
    s              *Storage
    updateLogger   *tglog.Logger
    logger         log.UnsafeInfoErrorLogger
    adminsNotifier notify.AutoNotifier
    superuserID    int64
}

type Settings struct {
    Token        string
    LastUpdateID int
    Storage      *Storage
    UpdateLogger *tglog.Logger
    Logger       log.UnsafeInfoErrorLogger
    SuperuserID  int64
}

func New(s *Settings) (*Bot, error) {
    b := Bot{
        s:            s.Storage,
        logger:       s.Logger,
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
                    b.logger.Error(errors.Wrap(err, "unable to log update"))
                }
            }()
            return true
        }),
        Client: &http.Client{
            Timeout: 0,
        },
        Reporter: func(err error) {
            b.logger.Error(err)
        },
    })
    if err != nil {
        return nil, err
    }
    b.tg = bot
    b.logger.Info(fmt.Sprintf("Authorized as @%s", b.tg.Me.Username))
    b.adminsNotifier = &notify.AdminsNotifier{
        ErrorLogger: s.Logger,
        Admins:      s.Storage.Admins,
        Notifiers: map[admins.NotifierType]notify.Notifier{
            admins.TelegramNotifier: tgnotify.NewNotifier(b.tg, notificationsCommand),
        },
    }
    if err := b.setupSuperuser(s.SuperuserID); err != nil {
        return nil, err
    }
    b.setupHandlers()
    return &b, nil
}

func (b *Bot) setupSuperuser(userID int64) error {
    userIDStr := strconv.FormatInt(userID, 10)
    if _, err := b.tg.ChatByID(userIDStr); err != nil {
        return errors.Wrapf(err,
            "can't get chat with superuser %q, check that superuser has a conversation with bot %q",
            userIDStr,
            b.tg.Me.Username)
    }
    _, err := b.s.Admins.CreateIfNotExists(context.Background(), &admins.Admin{
        ID:   userID,
        Role: admins.SuperuserRole,
        Notifications: &admins.Notifications{
            Status: []admins.NotifierType{
                admins.TelegramNotifier,
            },
        },
    })
    return err
}

func (b *Bot) setupHandlers() {
    b.tg.Handle(startCommand, handlers.MsgWithLog(b.logger,
        &handlers.Start{
            B: b.tg,
        }))
    b.tg.Handle(adminsCommand, handlers.MsgWithLog(b.logger,
        handlers.WithMsgFilters(
            handlers.NewAdmins(&handlers.AdminsSettings{
                Logger: b.logger,
                Tg:     b.tg,
                Storage: &handlers.AdminsStorage{
                    Admins: b.s.Admins,
                    Roles:  b.s.Roles,
                },
            }),
            filters.WithSender().IsAdminWithScopes(b.s.Admins,
                admins.AdminsReadScope,
            ),
        )))
    b.tg.Handle(notificationsCommand, handlers.MsgWithLog(b.logger,
        handlers.WithMsgFilters(handlers.NewNotifications(&handlers.NotificationsSettings{
            Logger: b.logger,
            Tg:     b.tg,
            Storage: &handlers.NotificationsStorage{
                Admins: b.s.Admins,
            },
        }),
            filters.WithSender().IsAdmin(b.s.Admins),
        )))
}

func (b *Bot) Start() {
    b.logger.Info("STARTED")
    if err := b.adminsNotifier.Notify(notify.StatusNotification(notify.StatusUp)); err != nil {
        b.logger.Error(err)
    }
    b.tg.Start()
}

func (b *Bot) Stop() {
    defer func() {
        if err := b.adminsNotifier.Notify(notify.StatusNotification(notify.StatusDown)); err != nil {
            b.logger.Error(err)
        }
    }()
    defer b.logger.Info("STOPPED")
    b.tg.Stop()
}
