package tg

import (
    "context"

    "github.com/mitinarseny/telego/bot/repo"
)

type User struct {
    repo.BaseModel `bson:",inline"`

    ID           int64   `bson:"_id,omitempty"`
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
