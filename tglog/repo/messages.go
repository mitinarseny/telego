package repo

import (
    "context"
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Message struct {
    createdAt *time.Time

    MessageID             int64
    From                  *User
    Date                  time.Time
    Chat                  *Chat
    ForwardFrom           *User
    ForwardFromChat       *Chat
    ForwardFromMessageID  *int64
    ForwardSignature      *string
    ForwardSenderName     *string
    ForwardDate           *time.Time
    ReplyToMessage        *Message
    EditDate              *time.Time
    MediaGroupID          *string
    AuthorSignature       *string
    Text                  *string
    Entities              []MessageEntity
    CaptionEntities       []MessageEntity
    Audio                 *Audio
    Document              *Document
    Animation             *Animation
    Game                  *Game
    Photo                 []PhotoSize
    Sticker               *Sticker
    Video                 *Video
    Voice                 *Voice
    VideoNote             *VideoNote
    Caption               *string
    Contact               *Contact
    Location              *Location
    Venue                 *Venue
    Poll                  *Poll
    NewChatMembers        []User
    LeftChatMember        *User
    NewChatTitle          *string
    NewChatPhoto          []PhotoSize
    DeleteChatPhoto       bool
    GroupChatCreated      bool
    SupergroupChatCreated bool
    ChannelChatCreated    bool
    MigrateToChatID       *int64
    MigrateFromChatID     *int64
    PinnedMessage         *Message
    Invoice               *Invoice
    SuccessfulPayment     *SuccessfulPayment
    ConnectedWebsite      *string
    PassportData          *PassportData
    ReplyMarkup           *InlineKeyboardMarkup
}

func (m *Message) ID() interface{} {
    return m.MessageID
}

func (m *Message) CreatedAt() *time.Time {
    return m.createdAt
}

type MessagesRepository interface {
    Create(ctx context.Context, messages ...*Message) ([]*Message, error)
}

func fromTelebotMessage(m *tb.Message) *Message {
    msg := new(Message)
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
    msg.Entities = make([]MessageEntity, 0, len(m.Entities))
    for _, e := range m.Entities {
        msg.Entities = append(msg.Entities, *fromTelebotMessageEntity(&e))
    }
    msg.CaptionEntities = make([]MessageEntity, 0, len(m.CaptionEntities))
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
        msg.Photo = append(make([]PhotoSize, 0, 1), *fromTelebotPhoto(m.Photo))
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
        msg.NewChatPhoto = append(make([]PhotoSize, 0, 1), *fromTelebotPhoto(m.NewGroupPhoto))
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
