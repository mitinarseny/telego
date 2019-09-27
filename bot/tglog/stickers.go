package tglog

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type MaskPositionPoint string

const (
    ForeheadMaskPositionPoint MaskPositionPoint = "forehead"
    EyesMaskPositionPoint     MaskPositionPoint = "eyes"
    MouthMaskPositionPoint    MaskPositionPoint = "mouth"
    ChinMaskPositionPoint     MaskPositionPoint = "chin"
)

type Sticker struct {
    FileID       string        `bson:"file_id,omitempty"`
    Width        int           `bson:"width,omitempty"`
    Height       int           `bson:"height,omitempty"`
    IsAnimated   bool          `bson:"is_animated,omitempty"`
    Thumb        *PhotoSize    `bson:"thumb,omitempty"`
    Emoji        *string       `bson:"emoji,omitempty"`
    SetName      *string       `bson:"set_name,omitempty"`
    MaskPosition *MaskPosition `bson:"mask_position,omitempty"`
    FileSize     *int64        `bson:"file_size,omitempty"`
}

type MaskPosition struct {
    Point  MaskPositionPoint `bson:"point,omitempty"`
    XShift float32           `bson:"x_shift,omitempty"`
    YShift float32           `bson:"y_shift,omitempty"`
    Scale  float32           `bson:"scale,omitempty"`
}

func fromTelebotSticker(s *tb.Sticker) *Sticker {
    st := new(Sticker)
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

func fromTelebotMaskPosition(m *tb.MaskPosition) *MaskPosition {
    return &MaskPosition{
        Point:  MaskPositionPoint(m.Feature),
        XShift: m.XShift,
        YShift: m.YShift,
        Scale:  m.Scale,
    }
}
