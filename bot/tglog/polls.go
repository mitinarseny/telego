package tglog

type Poll struct {
    ID       string       `bson:"id,omitempty"`
    Question string       `bson:"question,omitempty"`
    Options  []PollOption `bson:"options,omitempty"`
    IsClosed bool         `bson:"is_closed,omitempty"`
}

type PollOption struct {
    Text       string `bson:"text,omitempty"`
    VoterCount int    `bson:"voter_count,omitempty"`
}
