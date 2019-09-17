package bot

import (
    "github.com/mitinarseny/telego/administration/repo"
    "github.com/mitinarseny/telego/bot/filters"
    "github.com/mitinarseny/telego/bot/handlers"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MsgHandlerFunc func(*tb.Message) error

type MsgFilterFunc func(*tb.Message) (bool, error)

type Storage struct {
    Admins repo.AdminsRepo
    Roles  repo.RolesRepo
}

type Bot struct {
    handlers.Logger
    tg *tb.Bot
    s  *Storage
}

type Params struct {
    Logger  handlers.Logger
    Storage *Storage
}

func NewBot(bot *tb.Bot, params *Params) (*Bot, error) {
    b := &Bot{
        tg:     bot,
        s:      params.Storage,
        Logger: params.Logger,
    }

    bot.Handle("/start", handlers.MsgWithLog(b, &handlers.Start{
        B: bot,
    }))
    bot.Handle("/admins", handlers.MsgWithLog(b, handlers.WithMsgFilters(&handlers.Admins{
        Logger: b,
        B:      bot,
        Storage: &handlers.AdminsStorage{
            Admins: b.s.Admins,
        },
    }, filters.WithSender().IsAdminWithScopes(b.s.Admins, repo.AdminsReadScope))))
    return b, nil
}

func (b *Bot) Start() {
    b.tg.Start()
}

func (b *Bot) Stop() {
    b.tg.Stop()
}
