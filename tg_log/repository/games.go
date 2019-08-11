package repository

type Game struct {
    Title        string
    Description  string
    Photo        []PhotoSize
    Text         *string
    TextEntities []MessageEntity
    Animation    *Animation
}
