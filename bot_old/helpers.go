package bot_old

import "strings"

var mdEscaper = strings.NewReplacer("_", `\_`, "*", `\*`)

// escMd escapes all special Markdown symbols
func escMd(s string) string {
    return mdEscaper.Replace(s)
}
