package tg

type PhotoSize struct {
    FileID   string `bson:"file_id,omitempty"`
    Width    int    `bson:"width,omitempty"`
    Height   int    `bson:"height,omitempty"`
    FileSize *int64 `bson:"file_size,omitempty"`
}
