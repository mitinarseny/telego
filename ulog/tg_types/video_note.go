package tg_types

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type VideoNote struct {
    FileID   string
    Length   int
    Duration time.Duration
    Thumb    *PhotoSize
    FileSize *int64
}

func fromTelebotVideoNote(v *tb.VideoNote) *VideoNote {
    vn := new(VideoNote)
    vn.FileID = v.FileID
    vn.Length = v.Length
    vn.Duration = time.Duration(v.Duration)
    if v.Thumbnail != nil {
        vn.Thumb = fromTelebotPhoto(v.Thumbnail)
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vn.FileSize = &tmp
    }
    return vn
}
