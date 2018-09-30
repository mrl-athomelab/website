package template

import (
	tpl "html/template"
)

func AddFunction(a, b int) int {
	return a + b
}

func AsHTML(a string) tpl.HTML {
	return tpl.HTML(a)
}
