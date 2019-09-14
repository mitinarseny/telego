package utils

import "strings"

var mdEscaper = strings.NewReplacer("_", `\_`, "*", `\*`)

// EscMd escapes all special Markdown symbols
func EscMd(s string) string {
    return mdEscaper.Replace(s)
}
