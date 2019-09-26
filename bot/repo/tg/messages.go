package tg

import (
    "time"
)

type Message struct {
    MessageID             int64                 `bson:"_id,omitempty"`
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
