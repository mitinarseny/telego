package filters

import (
    tb "gopkg.in/tucnak/telebot.v2"
)

type and struct {
    msgFilters      []MsgFilter
    callbackFilters []CallbackFilter
}

func AndMsg(fs ...MsgFilter) *and {
    return &and{
        msgFilters: fs,
    }
}

func AndCallback(fs ...CallbackFilter) *and {
    return &and{
        callbackFilters: fs,
    }
}

func (f *and) FilterMsg(m *tb.Message) (bool, error) {
    for _, f := range f.msgFilters {
        switch passed, err := f.FilterMsg(m); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return true, nil
}

func (f *and) FilterCallback(c *tb.Callback) (bool, error) {
    for _, f := range f.callbackFilters {
        switch passed, err := f.FilterCallback(c); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return true, nil
}

type or struct {
    msgFilters      []MsgFilter
    callbackFilters []CallbackFilter
}

func OrMsg(fs ...MsgFilter) *or {
    return &or{
        msgFilters: fs,
    }
}

func OrCallback(fs ...CallbackFilter) *or {
    return &or{
        callbackFilters: fs,
    }
}

func (f *or) FilterMsg(m *tb.Message) (bool, error) {
    for _, f := range f.msgFilters {
        switch passed, err := f.FilterMsg(m); {
        case err != nil:
            return false, err
        case passed:
            return true, nil
        }
    }
    return false, nil
}

func (f *or) FilterCallback(c *tb.Callback) (bool, error) {
    for _, f := range f.callbackFilters {
        switch passed, err := f.FilterCallback(c); {
        case err != nil:
            return false, err
        case passed:
            return true, nil
        }
    }
    return false, nil
}
