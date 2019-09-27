package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Location struct {
    Longitude float32 `bson:"longitude,omitempty"`
    Latitude  float32 `bson:"latitude,omitempty"`
}

func fromTelebotLocation(l *tb.Location) *Location {
    return &Location{
        Longitude: l.Lng,
        Latitude:  l.Lat,
    }
}
