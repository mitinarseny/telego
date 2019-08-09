package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type InlineQuery struct {
    ID       string
    From     User
    Location *Location
    Query    string
    Offset   string
}

func fromTelebotQuery(q *tb.Query) *InlineQuery{
    iq := new(InlineQuery)
    iq.ID = q.ID
    iq.From = *fromTelebotUser(&q.From)
    if q.Location != nil {
        iq.Location = fromTelebotLocation(q.Location)
    }
    iq.Query = q.Text
    iq.Offset = q.Offset
    return iq
}
