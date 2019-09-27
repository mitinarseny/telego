package tglog

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Video struct {
    FileID   string        `bson:"file_id,omitempty"`
    Width    int           `bson:"width,omitempty"`
    Height   int           `bson:"height,omitempty"`
    Duration time.Duration `bson:"duration,omitempty"`
    Thumb    *PhotoSize    `bson:"thumb,omitempty"`
    MimeType *string       `bson:"mime_type,omitempty"`
    FileSize *int64        `bson:"file_size,omitempty"`
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
