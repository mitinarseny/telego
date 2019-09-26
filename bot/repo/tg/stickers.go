package tg

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
