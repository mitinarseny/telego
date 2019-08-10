package tg_log

import "github.com/mitinarseny/telego/tg_log/tg_types"

type UpdateLogger interface {
    LogUpdate(tg_types.Update) error
}
