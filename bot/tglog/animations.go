package tglog

import "time"

type Animation struct {
    FileID   string     `bson:"file_id,omitempty"`
    Width    int        `bson:"width,omitempty"`
    Height   int        `bson:"height,omitempty"`
    Duration time.Time  `bson:"duration,omitempty"`
    Thumb    *PhotoSize `bson:"thumb,omitempty"`
    FileName *string    `bson:"file_name,omitempty"`
    MimeType *string    `bson:"mime_type,omitempty"`
    FileSize *int64     `bson:"file_size,omitempty"`
}
