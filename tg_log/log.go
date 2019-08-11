package tg_log

import "github.com/mitinarseny/telego/tg_log/repository"

type UpdatesLogger interface {
    LogUpdate(*repository.Update) error
}

type BufferedUpdatesLogger interface {
    UpdatesLogger
    Close() error
}
