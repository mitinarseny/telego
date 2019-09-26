package tg

import (
    "time"
)

type Voice struct {
    FileID   string        `bson:"file_id,omitempty"`
    Duration time.Duration `bson:"duration,omitempty"`
    MimeType *string       `bson:"mime_type,omitempty"`
    FileSize *int64        `bson:"file_size,omitempty"`
}
