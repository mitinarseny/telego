package tglog

import (
    "time"

    tb "gopkg.in/tucnak/telebot.v2"
)

type Message struct {
    MessageID             int64                 `bson:"message_id,omitempty"`
    From                  *User                 `bson:"from,omitempty"`
    Date                  time.Time             `bson:"date,omitempty"`
    Chat                  *Chat                 `bson:"chat,omitempty"`
    ForwardFrom           *User                 `bson:"forward_from,omitempty"`
    ForwardFromChat       *Chat                 `bson:"forward_from_chat,omitempty"`
    ForwardFromMessageID  *int64                `bson:"forward_from_message_id,omitempty"`
    ForwardSignature      *string               `bson:"forward_signature,omitempty"`
    ForwardSenderName     *string               `bson:"forward_sender_name,omitempty"`
    ForwardDate           *time.Time            `bson:"forward_date,omitempty"`
    ReplyToMessage        *Message              `bson:"reply_to_message,omitempty"`
    EditDate              *time.Time            `bson:"edit_date,omitempty"`
    MediaGroupID          *string               `bson:"media_group_id,omitempty"`
    AuthorSignature       *string               `bson:"author_signature,omitempty"`
    Text                  *string               `bson:"text,omitempty"`
    Entities              []MessageEntity       `bson:"entities,omitempty"`
    CaptionEntities       []MessageEntity       `bson:"caption_entities,omitempty"`
    Audio                 *Audio                `bson:"audio,omitempty"`
    Document              *Document             `bson:"document,omitempty"`
    Animation             *Animation            `bson:"animation,omitempty"`
    Game                  *Game                 `bson:"game,omitempty"`
    Photo                 []PhotoSize           `bson:"photo,omitempty"`
    Sticker               *Sticker              `bson:"sticker,omitempty"`
    Video                 *Video                `bson:"video,omitempty"`
    Voice                 *Voice                `bson:"voice,omitempty"`
    VideoNote             *VideoNote            `bson:"video_note,omitempty"`
    Caption               *string               `bson:"caption,omitempty"`
    Contact               *Contact              `bson:"contact,omitempty"`
    Location              *Location             `bson:"location,omitempty"`
    Venue                 *Venue                `bson:"venue,omitempty"`
    Poll                  *Poll                 `bson:"poll,omitempty"`
    NewChatMembers        []User                `bson:"new_chat_members,omitempty"`
    LeftChatMember        *User                 `bson:"left_chat_member,omitempty"`
    NewChatTitle          *string               `bson:"new_chat_title,omitempty"`
    NewChatPhoto          []PhotoSize           `bson:"new_chat_photo,omitempty"`
    DeleteChatPhoto       bool                  `bson:"delete_chat_photo,omitempty"`
    GroupChatCreated      bool                  `bson:"group_chat_created,omitempty"`
    SupergroupChatCreated bool                  `bson:"supergroup_chat_created,omitempty"`
    ChannelChatCreated    bool                  `bson:"channel_chat_created,omitempty"`
    MigrateToChatID       *int64                `bson:"migrate_to_chat_id,omitempty"`
    MigrateFromChatID     *int64                `bson:"migrate_from_chat_id,omitempty"`
    PinnedMessage         *Message              `bson:"pinned_message,omitempty"`
    Invoice               *Invoice              `bson:"invoice,omitempty"`
    SuccessfulPayment     *SuccessfulPayment    `bson:"successful_payment,omitempty"`
    ConnectedWebsite      *string               `bson:"connected_website,omitempty"`
    PassportData          *PassportData         `bson:"passport_data,omitempty"`
    ReplyMarkup           *InlineKeyboardMarkup `bson:"reply_markup,omitempty"`
}

func fromTelebotMessage(m *tb.Message) *Message {
    msg := new(Message)
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
