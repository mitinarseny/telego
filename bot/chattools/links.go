package chattools

import (
    "github.com/pkg/errors"
)

func (b *bot) GetTelegramName(userID string) (string, error) {
    chat, err := b.b.ChatByID(userID)
    if err != nil {
        return "", errors.Wrapf(err, "can not get chat with %q", userID)
    }
    name := chat.FirstName
    if chat.LastName != "" {
        name += " " + chat.LastName
    }
    return name, nil
}

func (b *bot) UserLink(userID string) (string, error) {
    name, err := b.GetTelegramName(userID)
    if err != nil {
        return "", err
    }
    return "[" + name + "](tg://user?id=" + userID + ")", nil
}
