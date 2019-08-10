package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)
type ChosenInlineResult struct {
    ResultID        string
    From            User
    Location        *Location
    InlineMessageID *string
    Query           string
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
