package tg

type InlineQuery struct {
    ID       string    `bson:"_id,omitempty"`
    From     User      `bson:"from,omitempty"`
    Location *Location `bson:"location,omitempty"`
    Query    string    `bson:"query,omitempty"`
    Offset   string    `bson:"offset,omitempty"`
}
