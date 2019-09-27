package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Venue struct {
    Location       Location `bson:"location,omitempty"`
    Title          string   `bson:"title,omitempty"`
    Address        string   `bson:"address,omitempty"`
    FourSquareID   *string  `bson:"foursquare_id,omitempty"`
    FourSquareType *string  `bson:"foursquare_type,omitempty"`
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
