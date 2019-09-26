package tg

type Location struct {
    Longitude float32 `bson:"longitude,omitempty"`
    Latitude  float32 `bson:"latitude,omitempty"`
}
