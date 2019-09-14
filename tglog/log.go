package tglog

import "github.com/mitinarseny/telego/tglog/repo"

type UpdatesLogger interface {
    LogUpdates(...*repo.Update) error
}
