package tg

import (
    "time"
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
