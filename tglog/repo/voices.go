package repo

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Voice struct {
    FileID   string
    Duration time.Duration
    MimeType *string
    FileSize *int64
}

func fromTelebotVoice(v *tb.Voice) *Voice {
    vo := new(Voice)
    vo.FileID = v.FileID
    vo.Duration = time.Duration(v.Duration)
    if v.MIME != "" {
        vo.MimeType = &v.MIME
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vo.FileSize = &tmp
    }
    return vo
}
