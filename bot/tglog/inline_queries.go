package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type InlineQuery struct {
    ID       string    `bson:"id,omitempty"`
    From     User      `bson:"from,omitempty"`
    Location *Location `bson:"location,omitempty"`
    Query    string    `bson:"query,omitempty"`
    Offset   string    `bson:"offset,omitempty"`
}

func fromTelebotQuery(q *tb.Query) *InlineQuery {
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
