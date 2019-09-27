package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type PhotoSize struct {
    FileID   string `bson:"file_id,omitempty"`
    Width    int    `bson:"width,omitempty"`
    Height   int    `bson:"height,omitempty"`
    FileSize *int64 `bson:"file_size,omitempty"`
}

func fromTelebotPhoto(p *tb.Photo) *PhotoSize {
    ph := new(PhotoSize)
    ph.FileID = p.FileID
    ph.Width = p.Width
    ph.Height = p.Height
    if p.FileSize != 0 {
        tmp := int64(p.FileSize)
        ph.FileSize = &tmp
    }
    return ph
}
