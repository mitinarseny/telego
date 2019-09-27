package tglog

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type VideoNote struct {
    FileID   string        `bson:"file_id,omitempty"`
    Length   int           `bson:"length,omitempty"`
    Duration time.Duration `bson:"duration,omitempty"`
    Thumb    *PhotoSize    `bson:"thumb,omitempty"`
    FileSize *int64        `bson:"file_size,omitempty"`
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
