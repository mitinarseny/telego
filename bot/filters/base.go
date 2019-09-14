package filters

import (
    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

type MessageHandler func(*tb.Message) error

type Storage struct {
    Admins repo.AdminsRepo
    Roles repo.RolesRepo
}

type Filters struct {
    storage *Storage
}

func NewFilters(storage *Storage) *Filters {
    return &Filters{
        storage: storage,
    }
}
