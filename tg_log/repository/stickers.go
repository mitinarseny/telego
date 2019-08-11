package repository

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
    FileID       string
    Width        int
    Height       int
    IsAnimated   bool
    Thumb        *PhotoSize
    Emoji        *string
    SetName      *string
    MaskPosition *MaskPosition
    FileSize     *int64
}

type MaskPosition struct {
    Point  MaskPositionPoint
    XShift float32
    YShift float32
    Scale  float32
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
