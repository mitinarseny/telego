package repo

import (
    "time"

    "github.com/mitinarseny/telego/bot/repo/tg"
    tb "gopkg.in/tucnak/telebot.v2"
)

func fromTelebotAudio(a *tb.Audio) *tg.Audio {
    return &tg.Audio{
        FileID:    a.FileID,
        Duration:  time.Unix(int64(a.Duration), 0),
        Performer: nil,
        Title:     nil,
        MimeType:  nil,
        FileSize:  nil,
        Thumb:     nil,
    }
}

func fromTelebotChat(c *tb.Chat) *tg.Chat {
    return &tg.Chat{
        ID:        c.ID,
        Type:      tg.ChatType(c.Type),
        Title:     &c.Title,
        Username:  &c.Username,
        FirstName: &c.FirstName,
        LastName:  &c.LastName,
    }
}

func fromTelebotChosenInlineResult(r *tb.ChosenInlineResult) *tg.ChosenInlineResult {
    cr := new(tg.ChosenInlineResult)
    cr.ResultID = r.ResultID
    cr.From = *fromTelebotUser(&r.From)
    if r.Location != nil {
        cr.Location = fromTelebotLocation(r.Location)
    }
    if r.MessageID != "" {
        cr.InlineMessageID = &r.MessageID
    }
    cr.Query = r.Query
    return cr
}

func fromTelebotContact(c *tb.Contact) *tg.Contact {
    ct := new(tg.Contact)
    ct.PhoneNumber = c.PhoneNumber
    ct.FirstName = c.FirstName
    if c.LastName != "" {
        ct.LastName = &c.LastName
    }
    if c.UserID != 0 {
        tmp := int64(c.UserID)
        ct.UserID = &tmp
    }
    return ct
}

func fromTelebotDocument(d *tb.Document) *tg.Document {
    doc := new(tg.Document)
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

func fromTelebotQuery(q *tb.Query) *tg.InlineQuery {
    iq := new(tg.InlineQuery)
    iq.ID = q.ID
    iq.From = *fromTelebotUser(&q.From)
    if q.Location != nil {
        iq.Location = fromTelebotLocation(q.Location)
    }
    iq.Query = q.Text
    iq.Offset = q.Offset
    return iq
}

func fromTelebotLocation(l *tb.Location) *tg.Location {
    return &tg.Location{
        Longitude: l.Lng,
        Latitude:  l.Lat,
    }
}

func fromTelebotMessageEntity(e *tb.MessageEntity) *tg.MessageEntity {
    en := new(tg.MessageEntity)
    en.Type = tg.MessageEntityType(e.Type)
    en.Offset = e.Offset
    en.Length = e.Length
    if e.URL != "" {
        en.URL = &e.URL
    }
    if e.User != nil {
        en.User = fromTelebotUser(e.User)
    }
    return en
}

func fromTelebotMessage(m *tb.Message) *tg.Message {
    msg := new(tg.Message)
    if m == nil {
        return msg
    }
    msg.MessageID = int64(m.ID)
    if m.Sender != nil {
        msg.From = fromTelebotUser(m.Sender)
    }
    msg.Date = m.Time()
    if m.Chat != nil {
        msg.Chat = fromTelebotChat(m.Chat)
    }
    if m.ReplyTo != nil {
        msg.ReplyToMessage = fromTelebotMessage(m.ReplyTo)
    }
    if m.LastEdit != 0 {
        tmp := time.Unix(m.LastEdit, 0)
        msg.EditDate = &tmp
    }
    if m.Text != "" {
        msg.Text = &m.Text
    }
    msg.Entities = make([]tg.MessageEntity, 0, len(m.Entities))
    for _, e := range m.Entities {
        msg.Entities = append(msg.Entities, *fromTelebotMessageEntity(&e))
    }
    msg.CaptionEntities = make([]tg.MessageEntity, 0, len(m.CaptionEntities))
    for _, e := range m.CaptionEntities {
        msg.CaptionEntities = append(msg.CaptionEntities, *fromTelebotMessageEntity(&e))
    }
    if m.Audio != nil {
        msg.Audio = fromTelebotAudio(m.Audio)
    }
    if m.Document != nil {
        msg.Document = fromTelebotDocument(m.Document)
    }
    if m.Photo != nil {
        msg.Photo = append(make([]tg.PhotoSize, 0, 1), *fromTelebotPhoto(m.Photo))
    }
    if m.Sticker != nil {
        msg.Sticker = fromTelebotSticker(m.Sticker)
    }
    if m.Video != nil {
        msg.Video = fromTelebotVideo(m.Video)
    }
    if m.Voice != nil {
        msg.Voice = fromTelebotVoice(m.Voice)
    }
    if m.VideoNote != nil {
        msg.VideoNote = fromTelebotVideoNote(m.VideoNote)
    }
    if m.Caption != "" {
        msg.Caption = &m.Caption
    }
    if m.Contact != nil {
        msg.Contact = fromTelebotContact(m.Contact)
    }
    if m.Location != nil {
        msg.Location = fromTelebotLocation(m.Location)
    }
    if m.Venue != nil {
        msg.Venue = fromTelebotVenue(m.Venue)
    }
    if m.UserLeft != nil {
        msg.LeftChatMember = fromTelebotUser(m.UserLeft)
    }
    if m.NewGroupTitle != "" {
        msg.NewChatTitle = &m.NewGroupTitle
    }
    if m.NewGroupPhoto != nil {
        msg.NewChatPhoto = append(make([]tg.PhotoSize, 0, 1), *fromTelebotPhoto(m.NewGroupPhoto))
    }
    msg.DeleteChatPhoto = m.GroupPhotoDeleted
    msg.GroupChatCreated = m.GroupCreated
    msg.SupergroupChatCreated = m.SuperGroupCreated
    msg.ChannelChatCreated = m.ChannelCreated
    if m.MigrateTo != 0 {
        msg.MigrateToChatID = &m.MigrateTo
    }
    if m.MigrateFrom != 0 {
        msg.MigrateFromChatID = &m.MigrateFrom
    }
    if m.PinnedMessage != nil {
        msg.PinnedMessage = fromTelebotMessage(m.PinnedMessage)
    }
    return msg
}

func fromPreCheckoutQuery(q *tb.PreCheckoutQuery) *tg.PreCheckoutQuery {
    cq := new(tg.PreCheckoutQuery)
    cq.ID = q.ID
    cq.From = *fromTelebotUser(q.Sender)
    cq.Currency = q.Currency
    cq.TotalAmount = q.Total
    cq.InvoicePayload = q.Payload
    return cq
}

func fromTelebotPhoto(p *tb.Photo) *tg.PhotoSize {
    ph := new(tg.PhotoSize)
    ph.FileID = p.FileID
    ph.Width = p.Width
    ph.Height = p.Height
    if p.FileSize != 0 {
        tmp := int64(p.FileSize)
        ph.FileSize = &tmp
    }
    return ph
}

func fromTelebotSticker(s *tb.Sticker) *tg.Sticker {
    st := new(tg.Sticker)
    st.FileID = s.FileID
    st.Width = s.Width
    st.Height = s.Height
    if s.Thumbnail != nil {
        st.Thumb = fromTelebotPhoto(s.Thumbnail)
    }
    if s.Emoji != "" {
        st.Emoji = &s.Emoji
    }
    if s.SetName != "" {
        st.SetName = &s.SetName
    }
    if s.MaskPosition != nil {
        st.MaskPosition = fromTelebotMaskPosition(s.MaskPosition)
    }
    if s.FileSize != 0 {
        tmp := int64(s.FileSize)
        st.FileSize = &tmp
    }
    return st
}

func fromTelebotMaskPosition(m *tb.MaskPosition) *tg.MaskPosition {
    return &tg.MaskPosition{
        Point:  tg.MaskPositionPoint(m.Feature),
        XShift: m.XShift,
        YShift: m.YShift,
        Scale:  m.Scale,
    }
}

func fromTelebotUpdate(u *tb.Update) *tg.Update {
    upd := new(tg.Update)
    upd.UpdateID = int64(u.ID)
    if u.Message != nil {
        upd.Message = fromTelebotMessage(u.Message)
    }
    if u.EditedMessage != nil {
        upd.EditedMessage = fromTelebotMessage(u.EditedMessage)
    }
    if u.ChannelPost != nil {
        upd.ChannelPost = fromTelebotMessage(u.ChannelPost)
    }
    if u.EditedChannelPost != nil {
        upd.EditedChannelPost = fromTelebotMessage(u.EditedChannelPost)
    }
    if u.Query != nil {
        upd.InlineQuery = fromTelebotQuery(u.Query)
    }
    if u.ChosenInlineResult != nil {
        upd.ChosenInlineResult = fromTelebotChosenInlineResult(u.ChosenInlineResult)
    }
    if u.PreCheckoutQuery != nil {
        upd.PreCheckoutQuery = fromPreCheckoutQuery(u.PreCheckoutQuery)
    }
    return upd
}

func fromTelebotUser(u *tb.User) *tg.User {
    return &tg.User{
        ID:           int64(u.ID),
        IsBot:        false,
        FirstName:    u.FirstName,
        LastName:     &u.LastName,
        Username:     &u.Username,
        LanguageCode: &u.LanguageCode,
    }
}

func fromTelebotVenue(v *tb.Venue) *tg.Venue {
    ve := new(tg.Venue)
    ve.Location = *fromTelebotLocation(&v.Location)
    ve.Title = v.Title
    ve.Address = v.Address
    if v.FoursquareID != "" {
        ve.FourSquareID = &v.FoursquareID
    }
    if v.FoursquareType != "" {
        ve.FourSquareType = &v.FoursquareType
    }
    return ve
}

func fromTelebotVideoNote(v *tb.VideoNote) *tg.VideoNote {
    vn := new(tg.VideoNote)
    vn.FileID = v.FileID
    vn.Length = v.Length
    vn.Duration = time.Duration(v.Duration)
    if v.Thumbnail != nil {
        vn.Thumb = fromTelebotPhoto(v.Thumbnail)
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vn.FileSize = &tmp
    }
    return vn
}

func fromTelebotVideo(v *tb.Video) *tg.Video {
    vid := new(tg.Video)
    vid.FileID = v.FileID
    vid.Width = v.Width
    vid.Height = v.Height
    vid.Duration = time.Duration(v.Duration)
    if v.Thumbnail != nil {
        vid.Thumb = fromTelebotPhoto(v.Thumbnail)
    }
    if v.MIME != "" {
        vid.MimeType = &v.MIME
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vid.FileSize = &tmp
    }
    return vid
}

func fromTelebotVoice(v *tb.Voice) *tg.Voice {
    vo := new(tg.Voice)
    vo.FileID = v.FileID
    vo.Duration = time.Duration(v.Duration)
    if v.MIME != "" {
        vo.MimeType = &v.MIME
    }
    if v.FileSize != 0 {
        tmp := int64(v.FileSize)
        vo.FileSize = &tmp
    }
    return vo
}
