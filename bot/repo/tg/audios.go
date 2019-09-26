package tg

import (
    "time"
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
