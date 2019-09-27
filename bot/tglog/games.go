package tglog

type Game struct {
    Title        string          `bson:"title,omitempty"`
    Description  string          `bson:"description,omitempty"`
    Photo        []PhotoSize     `bson:"photo,omitempty"`
    Text         *string         `bson:"text,omitempty"`
    TextEntities []MessageEntity `bson:"text_entities,omitempty"`
    Animation    *Animation      `bson:"animation,omitempty"`
}
