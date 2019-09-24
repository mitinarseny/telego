package filters

import (
    "regexp"

    tb "gopkg.in/tucnak/telebot.v2"
)

type data struct {
    callbackParent CallbackFilter
    shouldMatch    *regexp.Regexp
}

func DataShouldMatch(re *regexp.Regexp) *data {
    return &data{
        shouldMatch: re,
    }
}

func (f *data) FilterCallback(c *tb.Callback) (bool, error) {
    if f.callbackParent != nil {
        switch passed, err := f.callbackParent.FilterCallback(c); {
        case err != nil:
            return false, err
        case !passed:
            return false, nil
        }
    }
    return f.shouldMatch.MatchString(c.Data), nil
}
