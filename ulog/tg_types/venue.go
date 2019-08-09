package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Venue struct {
    Location       Location
    Title          string
    Address        string
    FourSquareID   *string
    FourSquareType *string
}

func fromTelebotVenue(v *tb.Venue) *Venue {
    ve := new(Venue)
    ve.Location = *fromTelebotLocation(&v.Location)
    ve.Title = v.Title
    ve.Address = v.Address
    if v.FoursquareID != "" {
        ve.FourSquareID = &v.FoursquareID
    }
    if v.FoursquareType != "" {
        ve.FourSquareType = &v.FoursquareType
    }
    return ve
}
