package bot

import (
    "fmt"

    log "github.com/sirupsen/logrus"

    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgHandler func(*tb.Message) error

type MsgFilter func(*tb.Message) (bool, error)

type Storage struct {
    Admins repo.AdminsRepo
    Roles  repo.RolesRepo
}

type Bot struct {
    tg      *tb.Bot
    storage *Storage
}

func NewBot(bot *tb.Bot, storage *Storage) (*Bot, error) {
    b := &Bot{
        tg:      bot,
        storage: storage,
    }
    b.tg.Handle("/start", b.withLog(b.handleStart))
    b.tg.Handle("/help", b.withLog(b.handleHelp))
    b.tg.Handle("/roles", b.withLog(b.handleRoles))
    b.tg.Handle("/newrole", b.withLog(b.handleNewRole))
    b.tg.Handle("/admins", b.withLog(b.withFilters(b.handleAdmins, b.onlyAdminsWithScopes(repo.AdminsReadScope))))
    b.tg.Handle("/addadmin", b.withLog(b.withFilters(b.handleAddAdmin, b.onlyAdminsWithScopes(repo.AdminsScope))))
    b.tg.Handle("/stats", b.withLog(b.withFilters(b.handleStats, b.hasSender, b.onlyAdminsWithScopes(repo.StatsScope))))
    return b, nil
}

func (b *Bot) Start() {
    b.tg.Start()
}

func (b *Bot) Stop() {
    b.tg.Stop()
}

func (b *Bot) withLog(handler interface{}) interface{} {
    switch h := handler.(type) {
    case func(*tb.Message) error:
        return func(m *tb.Message) {
            if err := h(m); err != nil {
                log.WithFields(log.Fields{
                    "context": "BOT",
                }).Error(err)
            }
        }
    case MsgHandler:
        return func(m *tb.Message) {
            if err := h(m); err != nil {
                log.WithFields(log.Fields{
                    "context": "BOT",
                }).Error(err)
            }
        }
    default:
        panic(fmt.Sprintf("unknown handler type: %T", handler))
    }
}
func (b *Bot) withFilters(h MsgHandler, filters ...MsgFilter) MsgHandler {
    return func(m *tb.Message) error {
        for _, f := range filters {
            passed, err := f(m)
            if err != nil {
                return err
            }
            if !passed {
                return nil
            }
        }
        return h(m)
    }
}
