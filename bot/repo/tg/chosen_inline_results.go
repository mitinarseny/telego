package tg

type ChosenInlineResult struct {
    ResultID        string    `bson:"result_id,omitempty"`
    From            User      `bson:"from,omitempty"`
    Location        *Location `bson:"location,omitempty"`
    InlineMessageID *string   `bson:"inline_message_id,omitempty"`
    Query           string    `bson:"query,omitempty"`
}
