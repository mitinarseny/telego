package tg_types

import "time"

type Animation struct {
    FileID   string
    Width    int
    Height   int
    Duration time.Time
    Thumb    *PhotoSize
    FileName *string
    MimeType *string
    FileSize *int64
}
