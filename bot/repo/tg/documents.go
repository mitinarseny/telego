package tg

type Document struct {
    FileID   string     `bson:"file_id,omitempty"`
    Thumb    *PhotoSize `bson:"thumb,omitempty"`
    FileName *string    `bson:"file_name,omitempty"`
    MimeType *string    `bson:"mime_type,omitempty"`
    FileSize *int64     `bson:"file_size,omitempty"`
}
