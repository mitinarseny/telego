package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Contact struct {
    PhoneNumber string
    FirstName   string
    LastName    *string
    UserID      *int64
    VCard       *string
}

func fromTelebotContact(c *tb.Contact) *Contact {
    ct := new(Contact)
    ct.PhoneNumber = c.PhoneNumber
    ct.FirstName = c.FirstName
    if c.LastName != "" {
        ct.LastName = &c.LastName
    }
    if c.UserID != 0 {
        tmp := int64(c.UserID)
        ct.UserID = &tmp
    }
    return ct
}
