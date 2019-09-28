package tglog

import (
    "context"

    tb "gopkg.in/tucnak/telebot.v2"
)

type User struct {
    ID           int64   `bson:"id,omitempty"`
    IsBot        bool    `bson:"is_bot,omitempty"`
    FirstName    string  `bson:"first_name,omitempty"`
    LastName     *string `bson:"last_name,omitempty"`
    Username     *string `bson:"username,omitempty"`
    LanguageCode *string `bson:"language_code,omitempty"`
}

type UsersRepo interface {
    Create(ctx context.Context, users ...*User) ([]*User, error)
    CreateIfNotExist(ctx context.Context, users ...*User) ([]*User, error)
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
