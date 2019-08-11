package repository

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type User struct {
    ID           int64
    IsBot        bool
    FirstName    string
    LastName     *string
    Username     *string
    LanguageCode *string
}

func fromTelebotUser(u *tb.User) *User {
    return &User{
        ID:           int64(u.ID),
        IsBot:        false,
        FirstName:    u.FirstName,
        LastName:     &u.LastName,
        Username:     &u.Username,
        LanguageCode: &u.LanguageCode,
    }
}
