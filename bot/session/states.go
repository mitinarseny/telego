package session

type State map[string]interface{}

type MessageStateRepo interface {
    Get(chatID, messageID int64, keys ...string) (State, error)
    Replace(chatID, messageID int64, s State) error
    Update(chatID, messageID int64, upd State) error
}
