package tg_types

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type Document struct {
    FileID   string
    Thumb    *PhotoSize
    FileName *string
    MimeType *string
    FileSize *int64
}

func fromTelebotDocument( d *tb.Document) *Document {
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
