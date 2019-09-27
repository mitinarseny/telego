package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type ChosenInlineResult struct {
    ResultID        string    `bson:"result_id,omitempty"`
    From            User      `bson:"from,omitempty"`
    Location        *Location `bson:"location,omitempty"`
    InlineMessageID *string   `bson:"inline_message_id,omitempty"`
    Query           string    `bson:"query,omitempty"`
}

func fromTelebotChosenInlineResult(r *tb.ChosenInlineResult) *ChosenInlineResult {
    cr := new(ChosenInlineResult)
    cr.ResultID = r.ResultID
    cr.From = *fromTelebotUser(&r.From)
    if r.Location != nil {
        cr.Location = fromTelebotLocation(r.Location)
    }
    if r.MessageID != "" {
        cr.InlineMessageID = &r.MessageID
    }
    cr.Query = r.Query
    return cr
}
