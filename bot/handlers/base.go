package handlers

import (
    "github.com/mitinarseny/telego/administration/repo"
    tb "gopkg.in/tucnak/telebot.v2"
)

type Storage struct {
    Admins repo.AdminsRepo
    Roles repo.RolesRepo
}

type Handler struct {
    Bot *tb.Bot
    Storage *Storage
}
