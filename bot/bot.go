package bot

import (
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
    tg *tb.Bot
    storage *Storage
}

func NewBot(bot *tb.Bot, storage *Storage) (*Bot, error) {
    b := &Bot{
        tg:     bot,
        storage: storage,
    }
    b.tg.Handle("/start", b.withLogAndFilters(b.handleStart))
    b.tg.Handle("/stats", b.withLogAndFilters(b.handleStats, b.onlyAdminsWithScopes(repo.StatsScope)))
    b.tg.Handle("/admins", b.withLogAndFilters(b.handleAdmins, b.onlyAdminsWithScopes(repo.AdminsReadScope)))
    b.tg.Handle("/addadmin", b.withLogAndFilters(b.handleAddAdmin, b.onlyAdminsWithScopes(repo.AdminsScope)))
    return b, nil
}

func (b *Bot) Start() {
    b.tg.Start()
}

func (b *Bot) Stop() {
    b.tg.Stop()
}

func (b *Bot) withLogAndFilters(h MsgHandler, filters ...MsgFilter) func(*tb.Message) {
    return func(m *tb.Message) {
        for _, f := range filters {
            passed, err := f(m)
            if err != nil {
                log.WithFields(log.Fields{
                    "context": "BOT",
                    "filter": f,
                }).Error(err)
                return
            }
            if !passed {
                return
            }
        }
        if err := h(m); err != nil {
            log.WithFields(log.Fields{
                "context": "BOT",
                "handler": h,
            }).Error(err)
        }
    }
}
