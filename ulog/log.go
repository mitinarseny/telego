package ulog

import "github.com/mitinarseny/telego/ulog/tg_types"

type UpdateLogger interface {
    LogUpdate(tg_types.Update) error
}
