package repo

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Location struct {
    Longitude float32
    Latitude  float32
}

func fromTelebotLocation(l *tb.Location) *Location {
    return &Location{
        Longitude: l.Lng,
        Latitude:  l.Lat,
    }
}
