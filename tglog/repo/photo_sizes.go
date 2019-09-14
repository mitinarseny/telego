package repo

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type PhotoSize struct {
    FileID   string
    Width    int
    Height   int
    FileSize *int64
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
