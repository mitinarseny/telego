package tg

import (
    "time"
)

type VideoNote struct {
    FileID   string        `bson:"file_id,omitempty"`
    Length   int           `bson:"length,omitempty"`
    Duration time.Duration `bson:"duration,omitempty"`
    Thumb    *PhotoSize    `bson:"thumb,omitempty"`
    FileSize *int64        `bson:"file_size,omitempty"`
}
