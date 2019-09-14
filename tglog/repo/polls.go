package repo

type Poll struct {
    ID       string
    Question string
    Options  []PollOption
    IsClosed bool
}

type PollOption struct {
    Text       string
    VoterCount int
}