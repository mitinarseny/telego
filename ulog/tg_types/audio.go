package tg_types

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Audio struct {
    FileID    string
    Duration  time.Time
    Performer *string
    Title     *string
    MimeType  *string
    FileSize  *int64
    Thumb     *PhotoSize
}

func fromTelebotAudio(a *tb.Audio) *Audio {
    return &Audio{
        FileID:    a.FileID,
        Duration:  time.Unix(int64(a.Duration), 0),
        Performer: nil,
        Title:     nil,
        MimeType:  nil,
        FileSize:  nil,
        Thumb:     nil,
    }
}
