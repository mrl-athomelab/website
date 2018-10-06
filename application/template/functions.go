package template

import (
	tpl "html/template"
	"time"
)

func AddFunction(a, b int) int {
	return a + b
}

func AsHTML(a string) tpl.HTML {
	return tpl.HTML(a)
}

func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}
