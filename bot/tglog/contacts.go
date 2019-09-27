package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Contact struct {
    PhoneNumber string  `bson:"phone_number,omitempty"`
    FirstName   string  `bson:"first_name,omitempty"`
    LastName    *string `bson:"last_name,omitempty"`
    UserID      *int64  `bson:"user_id,omitempty"`
    VCard       *string `bson:"vcard,omitempty"`
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
