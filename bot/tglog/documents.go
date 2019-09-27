package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Document struct {
    FileID   string     `bson:"file_id,omitempty"`
    Thumb    *PhotoSize `bson:"thumb,omitempty"`
    FileName *string    `bson:"file_name,omitempty"`
    MimeType *string    `bson:"mime_type,omitempty"`
    FileSize *int64     `bson:"file_size,omitempty"`
}

func fromTelebotDocument(d *tb.Document) *Document {
    doc := new(Document)
    doc.FileID = d.FileID
    doc.Thumb = fromTelebotPhoto(d.Thumbnail)
    if d.FileName != "" {
        doc.FileName = &d.FileName
    }
    if d.MIME != "" {
        doc.MimeType = &d.MIME
    }
    if d.FileSize != 0 {
        tmp := int64(d.FileSize)
        doc.FileSize = &tmp
    }
    return doc
}
