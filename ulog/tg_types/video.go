package tg_types

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Video struct {
    FileID   string
    Width    int
    Height   int
    Duration time.Duration
    Thumb    *PhotoSize
    MimeType *string
    FileSize *int64
}

func fromTelebotVideo(v *tb.Video) *Video {
    vid := new(Video)
    vid.FileID = v.FileID
    vid.Width = v.Width
    vid.Height = v.Height
    vid.Duration = time.Duration(v.Duration)
    if v.Thumbnail != nil {
        vid.Thumb = fromTelebotPhoto(v.Thumbnail)
    }
    if v.MIME != "" {
        vid.MimeType = &v.MIME
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vid.FileSize = &tmp
    }
    return vid
}
