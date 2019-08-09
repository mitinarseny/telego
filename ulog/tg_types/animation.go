package tg_types

type Animation struct {
    FileID   string
    Width    int
    Height   int
    Duration int
    Thumb    *PhotoSize
    FileName *string
    MimeType *string
    FileSize *int64
}
