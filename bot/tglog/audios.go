package tglog

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Audio struct {
    FileID    string     `bson:"file_id,omitempty"`
    Duration  time.Time  `bson:"duration,omitempty"`
    Performer *string    `bson:"performer,omitempty"`
    Title     *string    `bson:"title,omitempty"`
    MimeType  *string    `bson:"mime_type,omitempty"`
    FileSize  *int64     `bson:"file_size,omitempty"`
    Thumb     *PhotoSize `bson:"thumb,omitempty"`
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
